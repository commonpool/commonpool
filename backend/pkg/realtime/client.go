package realtime

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/commonpool/backend/pkg/handler"
	"github.com/commonpool/backend/pkg/mq"
	"github.com/gorilla/websocket"
	"github.com/labstack/gommon/log"
	"go.uber.org/zap"
	"time"
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

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

type Client struct {
	hub                 *Hub
	websocketConnection *websocket.Conn
	queue               mq.Queue
	send                chan []byte
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.websocketConnection.Close()
		c.queue.Close()
		log.Print("closing readPump")
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

	ctx, l := handler.GetCtx(ctx, "eventPump")

	err := c.queue.Consume(ctx, mq.NewConsumerConfig(func(delivery mq.Delivery) error {

		l.Debug("received message from RabbitMQ")

		err := delivery.Acknowledger.Ack(delivery.DeliveryTag, false)
		if err != nil {
			l.Error("could not acknowledge delivery", zap.Error(err))
			return err
		}

		var event mq.Event
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

		return nil

	}))

	if err != nil {
		return err
	}

	select {
	case <-c.queue.Closed():
		l.Info("event pump stopping")
		return nil
	}

}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.websocketConnection.Close()
		log.Print("closing writePump")
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
