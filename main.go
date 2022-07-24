package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"net/http"
	"tictactoe/api"
	"tictactoe/pkg/websocket"
)

func main() {
	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
	}))
	r.Use(middleware.Logger)
	r.Get("/", api.GetBoard)
	r.Post("/move", api.ReceiveMove)
	r.HandleFunc("/ws", websocket.ServeWs)
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		return
	}
}
