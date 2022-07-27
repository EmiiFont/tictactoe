package websocket

import (
	"fmt"
	"github.com/google/uuid"
)

type Pool struct {
	Register   chan *Client
	Unregister chan *Client
	Clients    map[uuid.UUID][]*Client
	SendToRoom chan Message
}

func NewPool() *Pool {
	return &Pool{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[uuid.UUID][]*Client),
		SendToRoom: make(chan Message),
	}
}

func (pool *Pool) Start() {
	for {
		select {
		case client := <-pool.Register:
			pool.Clients[client.RoomId] = append(pool.Clients[client.RoomId], client)
			fmt.Println("Size of Connection Pool: ", len(pool.Clients[client.RoomId]))
			for _, client := range pool.Clients[client.RoomId] {
				fmt.Println(client)
				fmt.Println(client.RoomId)
				err := client.Conn.WriteJSON(Message{Type: 1, Body: "New User Joined..."})
				if err != nil {
					return
				}
			}
			break
		case client := <-pool.Unregister:
			delete(pool.Clients, client.RoomId)
			fmt.Println("Size of Connection Pool: ", len(pool.Clients))
			theClientInRoom := pool.Clients[client.RoomId]
			for _, client := range theClientInRoom {
				err := client.Conn.WriteJSON(Message{Type: 1, Body: "User Disconnected..."})
				if err != nil {
					return
				}
			}
			break
		case message := <-pool.SendToRoom:
			fmt.Printf("Sending message to all clients in Room %s\n", message.RoomId)
			for _, client := range pool.Clients[message.RoomId] {
				if err := client.Conn.WriteJSON(message); err != nil {
					fmt.Println(err)
					return
				}
			}
		}
	}
}
