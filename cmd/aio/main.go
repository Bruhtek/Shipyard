package main

import (
	"Shipyard/internal/api/actions"
	"Shipyard/internal/api/env"
	"Shipyard/internal/api/websocket"
	"Shipyard/internal/intervals"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"time"
)

func main() {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Compress(5,
		"text/html",
		"text/css",
		"application/json",
		"application/javascript"))
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	intervals.SetupIntervals()

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
