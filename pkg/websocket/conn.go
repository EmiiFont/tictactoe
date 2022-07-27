package websocket

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

// We'll need to define an Upgrader
// this will require a Read and Write buffer size
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,

	// We'll need to check the origin of our connection
	// this will allow us to make requests from our React
	// development server to here.
	// For now, we'll do no checking and just allow any connection
	CheckOrigin: func(r *http.Request) bool { return true },
}

// ServeWs define our WebSocket endpoint
func ServeWs(pool *Pool, w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Host)
	roomId := chi.URLParam(r, "roomId")
	// upgrade this connection to a WebSocket
	// connection
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	udi, _ := uuid.Parse(roomId)
	client := &Client{
		Conn:   ws,
		Pool:   pool,
		RoomId: udi,
	}

	pool.Register <- client
	// listen indefinitely for new messages coming
	// through on our WebSocket connection
	client.Read()
}
