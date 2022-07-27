package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"net/http"
	"tictactoe/api"
	"tictactoe/pkg/websocket"
)

func setupRoutes(chiRouter *chi.Mux) {
	pool := websocket.NewPool()
	go pool.Start()

	chiRouter.HandleFunc("/ws/{roomId}", func(w http.ResponseWriter, r *http.Request) {
		websocket.ServeWs(pool, w, r)
	})
}

func main() {
	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
	}))
	r.Use(middleware.Logger)
	r.Get("/", api.GetBoard)
	r.Post("/move", api.ReceiveMove)
	r.Get("/newGame", api.NewGame)
	setupRoutes(r)
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		return
	}
}
