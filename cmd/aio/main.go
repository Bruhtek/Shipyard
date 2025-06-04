package main

import (
	"Shipyard/database"
	"Shipyard/internal/api/actions"
	"Shipyard/internal/api/env"
	"Shipyard/internal/api/websocket"
	"Shipyard/internal/env_manager"
	"Shipyard/internal/intervals"
	"Shipyard/internal/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"time"
)

func main() {
	r := chi.NewRouter()

	logger.Init(os.Getenv("ENV") == "development")
	r.Use(logger.HttpLogger)
	
	database.InitializeDatabase()
	env_manager.InitializeEnvManager()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://localhost:*", "http://localhost:*"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.Use(middleware.RequestID)
	r.Use(middleware.Compress(5,
		"text/html",
		"text/css",
		"application/json",
		"application/javascript"))
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

	log.Info().Int("port", 4000).Msg("Starting server")
	http.ListenAndServe(":4000", r)
}
