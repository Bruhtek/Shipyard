package main

import (
	"Shipyard/internal/api/actions"
	"Shipyard/internal/api/env"
	"Shipyard/internal/env_manager"
	"Shipyard/internal/intervals"
	"Shipyard/internal/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"time"
)

func main() {
	r := chi.NewRouter()

	logger.Init(os.Getenv("ENV") == "development")
	r.Use(logger.HttpLogger)

	env_manager.InitializeEnvManager(true)

	r.Use(middleware.RequestID)
	r.Use(middleware.Compress(5,
		"text/html",
		"text/css",
		"application/json",
		"application/javascript"))
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	intervals.SetupIntervals()
	intervals.SetupHeartbeat()

	envRouter := env.GetEnvRouter()
	actionsRouter := actions.GetActionsRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	r.Mount("/api/env", envRouter)
	r.Mount("/api/actions", actionsRouter)

	log.Info().Int("port", 4333).Msg("Starting server")
	http.ListenAndServe(":4333", r)
}
