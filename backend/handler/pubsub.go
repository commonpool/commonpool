package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/commonpool/backend/amqp"
	"github.com/commonpool/backend/auth"
	"github.com/commonpool/backend/model"
	"github.com/commonpool/backend/utils"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"log"
	"net/http"
	"time"
)

type Client struct {
	hub                 *Hub
	websocketConnection *websocket.Conn
	amqpChannel         amqp.AmqpChannel
	send                chan []byte
	id                  string
	userKey             model.UserKey
	queueName           string
	consumerKey         string
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.websocketConnection.Close()
	}()
	c.websocketConnection.SetReadLimit(maxMessageSize)
	c.websocketConnection.SetReadDeadline(time.Now().Add(pongWait))
	c.websocketConnection.SetPongHandler(func(string) error { c.websocketConnection.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.websocketConnection.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		c.hub.broadcast <- message
	}
}

func (c *Client) eventPump(ctx context.Context) error {

	ctx, l := GetCtx(ctx, "eventPump")

	l.Debug("event pump")

	ch, err := c.amqpChannel.Consume(ctx, c.queueName, c.consumerKey, false, false, false, false, nil)
	if err != nil {
		l.Error("could not consume amqp channel", zap.Error(err))
		return err
	}

	for delivery := range ch {

		l.Debug("received message from RabbitMQ")

		err := delivery.Acknowledger.Ack(delivery.DeliveryTag, false)
		if err != nil {
			l.Error("could not acknowledge delivery", zap.Error(err))
			return err
		}

		var event amqp.Event
		err = json.Unmarshal(delivery.Body, &event)
		if err != nil {
			l.Error("could not unmarshal event", zap.Error(err))
			return err
		}

		js, err := json.Marshal(event)
		if err != nil {
			l.Error("could not marshal event", zap.Error(err))
			return err
		}

		c.send <- js

	}

	return nil

}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.websocketConnection.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.websocketConnection.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.websocketConnection.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.websocketConnection.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.websocketConnection.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.websocketConnection.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				_ = client.amqpChannel.Close()
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

func NewClient(hub *Hub, conn *websocket.Conn, amqpChannel amqp.AmqpChannel, queueName string, key string) *Client {

	return &Client{
		hub:                 hub,
		websocketConnection: conn,
		send:                make(chan []byte, 256),
		amqpChannel:         amqpChannel,
		queueName:           queueName,
		consumerKey:         key,
	}
}

func (h *Handler) Websocket(c echo.Context) error {

	ctx, l := GetEchoContext(c, "Websocket")

	l.Debug("initializing websocket connection")

	l.Debug("getting user session")

	userSession, err := auth.GetUserSession(ctx)
	if err != nil {
		l.Error("could not get user session", zap.Error(err))
		return err
	}
	userKey := userSession.GetUserKey()

	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	l.Debug("upgrading websocket connection")

	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		l.Error("could not upgrade websocket connection", zap.Error(err))
		return err
	}
	defer ws.Close()

	l.Debug("creating new hub")

	hub := newHub()
	go hub.run()

	consumerKey := utils.ShortUuid(uuid.NewV4())
	queueName := "chat.ws." + userKey.String() + "." + consumerKey

	l.Info("registering hub client", zap.String("queue_name", queueName), zap.String("consumer_key", consumerKey))

	amqpChannel, err := h.amqp.GetChannel()
	if err != nil {
		l.Error("cold not get amqp channel", zap.Error(err))
		return err
	}
	defer amqpChannel.Close()

	l.Debug("creating user exchange")

	userExchangeName, err := h.chatService.CreateUserExchange(ctx, userKey)
	if err != nil {
		l.Error("could not create user exchange", zap.Error(err))
		return err
	}

	l.Debug("declaring queue for websocket connection", zap.String("queue_name", queueName))

	err = amqpChannel.QueueDeclare(ctx, queueName, false, true, false, false, nil)
	if err != nil {
		l.Error("could not declare websocket amqp queue", zap.Error(err))
		return err
	}

	l.Debug("binding websocket consumer queue to user exchange")

	err = amqpChannel.QueueBind(ctx, queueName, "", userExchangeName, false, nil)
	if err != nil {
		l.Error("could not bind consumer queue to exchange", zap.Error(err))
		return err
	}

	client := NewClient(hub, ws, amqpChannel, queueName, consumerKey)

	c.Logger().Info("registering hub")
	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()

	return client.eventPump(ctx)

	// ch, err := h.amqp.RegisterWsChannel(ctx, clientId, userKey)
	// if err != nil {
	// 	l.Error("could not register ws amqpChannel", zap.Error(err))
	// 	return err
	// }
	//
	// err := h.groupService.RegisterUserAmqpSubscriptions(ctx)
	//
	// c.Logger().Info("getting subscriptions")
	// getSubscriptions := chat.NewGetSubscriptions(userKey, 100, 0)
	// subs, err := h.chatStore.GetSubscriptionsForUser(ctx, getSubscriptions)
	// if err != nil {
	// 	l.Error("could not get subscriptions", zap.Error(err))
	// 	return err
	// }
	//
	// c.Logger().Info("creating subscriptions for user")
	// for _, item := range subs.Subscriptions.Items {
	// 	err := h.amqp.BindUserExchangeToChannel(ctx, item.GetChannelKey(), userKey)
	// 	if err != nil {
	// 		l.Error("could not create rabbitmq binding for subscription",
	// 			zap.String("channel_id", item.ChannelID),
	// 			zap.String("user_id", item.UserID),
	// 			zap.Error(err),
	// 		)
	// 		return err
	// 	}
	// }
	//
	// membershipStatus := group.ApprovedMembershipStatus
	// getMemberships := group.NewGetMembershipsForUserRequest(userKey, &membershipStatus)
	// getMembershipsResponse := h.groupStore.GetMembershipsForUser(getMemberships)
	// if getMembershipsResponse.Error != nil {
	// 	l.Error("could not get group memberships for user", zap.Error(err))
	// 	return getMembershipsResponse.Error
	// }
	//
	// for _, membership := range getMembershipsResponse.Memberships.Items {
	// 	err = h.amqp.RegisterUserMembershipBinding(membership.GetUserKey(), membership.GetGroupKey())
	// 	if err != nil {
	// 		l.Error("could not register user membership binding", zap.Error(err))
	// 		return err
	// 	}
	// }

}
