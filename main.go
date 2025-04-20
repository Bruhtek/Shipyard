package main

import (
	"Shipyard/api/actions"
	"Shipyard/api/env"
	"Shipyard/api/websocket"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"time"
)

func main() {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	envRouter := env.GetEnvRouter()
	actionsRouter := actions.GetActionsRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})
	r.Get("/ws", websocket.HandleWebsocketConnection)

	r.Mount("/api/env", envRouter)
	r.Mount("/api/actions", actionsRouter)

	println("Starting server on port 4000")
	http.ListenAndServe(":4000", r)
}
