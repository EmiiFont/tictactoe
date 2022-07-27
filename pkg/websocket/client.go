package websocket

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"log"
)

type Client struct {
	ID     string
	Conn   *websocket.Conn
	Pool   *Pool
	RoomId uuid.UUID
}

type Message struct {
	Type   int       `json:"type"`
	Body   string    `json:"body"`
	RoomId uuid.UUID `json:"roomId"`
}

func (c *Client) Read() {
	defer func() {
		c.Pool.Unregister <- c
		c.Conn.Close()
	}()

	for {
		messageType, p, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		message := Message{Type: messageType, Body: string(p), RoomId: c.RoomId}
		c.Pool.SendToRoom <- message
		fmt.Printf("Message Received: %+v\n", message)
	}
}
