package realtime

import (
	"github.com/commonpool/backend/pkg/mq"
	"github.com/gorilla/websocket"
)

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
				if client.amqpChannel != nil {
					_ = client.amqpChannel.Close()
				}
			}
		case _ = <-h.broadcast:

		}
	}
}

func NewClient(hub *Hub, conn *websocket.Conn, amqpChannel mq.Channel, queueName *string, key *string) *Client {
	return &Client{
		hub:                 hub,
		websocketConnection: conn,
		send:                make(chan []byte, 256),
		amqpChannel:         amqpChannel,
		queueName:           queueName,
		consumerKey:         key,
	}
}

func NewAnonymousClient(hub *Hub, conn *websocket.Conn) *Client {
	return &Client{
		hub:                 hub,
		websocketConnection: conn,
		send:                make(chan []byte, 256),
		amqpChannel:         nil,
		queueName:           nil,
		consumerKey:         nil,
	}
}
