package remote

import (
	"Shipyard/internal/env_manager"
	"Shipyard/internal/remote_environment"
	"context"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func GetRemoteMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := chi.URLParam(r, "remote_key")

		remote := env_manager.EnvManager.GetEnvByKey(key)

		if remote == nil || remote.GetEnvDescription().EnvType != "remote" {
			http.Error(w, "Remote not found", http.StatusNotFound)
			return
		}
		ctx := context.WithValue(r.Context(), "remote", remote)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetRemoteRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Route("/{remote_key}", func(r chi.Router) {
		r.Use(GetRemoteMiddleware)

		r.Get("/heartbeat", func(w http.ResponseWriter, r *http.Request) {
			remote := r.Context().Value("remote").(*remote_environment.RemoteEnvironment)

			remote.Heartbeat()
			requested := remote.IsRequested()

			if requested {
				w.WriteHeader(http.StatusAccepted)
				w.Write([]byte("CONNECT"))
				return
			}

			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		})

		r.Get("/ws", remote_environment.HandleRemoteWebsocketConnection)
	})

	return r
}
