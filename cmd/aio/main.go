package main

import (
	"Shipyard"
	"Shipyard/database"
	"Shipyard/internal/api/actions"
	"Shipyard/internal/api/env"
	"Shipyard/internal/api/remote"
	"Shipyard/internal/api/websocket"
	"Shipyard/internal/env_manager"
	"Shipyard/internal/intervals"
	"Shipyard/internal/logger"
	"embed"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/rs/zerolog/log"
	"io/fs"
	"net/http"
	"os"
	"time"
)

var staticWebContent embed.FS = Shipyard.WebContent

func main() {
	r := chi.NewRouter()

	logger.Init(os.Getenv("ENV") == "development")
	r.Use(logger.HttpLogger)

	database.InitializeDatabase()
	env_manager.InitializeEnvManager(false)

	if os.Getenv("ENV") == "development" {
		r.Use(cors.Handler(cors.Options{
			AllowedOrigins:   []string{"https://localhost:*", "http://localhost:*"},
			AllowCredentials: true,
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			MaxAge:           300, // Maximum value not ignored by any of major browsers
		}))
	}

	r.Use(middleware.RequestID)
	r.Use(middleware.Compress(5,
		"text/html",
		"text/css",
		"application/json",
		"application/javascript"))
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	intervals.SetupIntervals(false)

	envRouter := env.GetEnvRouter()
	actionsRouter := actions.GetActionsRouter()
	remoteRouter := remote.GetRemoteRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})
	r.Get("/ws", websocket.HandleWebsocketConnection)

	r.Mount("/api/env", envRouter)
	r.Mount("/api/actions", actionsRouter)
	r.Mount("/api/remote", remoteRouter)

	content, err := fs.Sub(staticWebContent, "web/build")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load embedded web files")
	}

	staticServer := http.FileServer(http.FS(content))
	r.Handle("/", staticServer)
	r.Handle("/favicon.ico", staticServer)
	r.Handle("/_app/*", staticServer)

	r.HandleFunc("/*", func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = "/" // Replace the request path
		staticServer.ServeHTTP(w, r)
	})
	log.Info().Int("port", 4000).Msg("Starting server")
	http.ListenAndServe(":4000", r)
}
