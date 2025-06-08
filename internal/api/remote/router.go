package remote

import (
	"Shipyard/internal/env_manager"
	"context"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func RemoteMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := chi.URLParam(r, "key")

		if env := env_manager.EnvManager.GetRemoteEnv(key); env != nil {
			ctx := context.WithValue(r.Context(), "env", env)

			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
	})
}

func GetRemoteRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Route("/{key}", func(r chi.Router) {
		r.Use(RemoteMiddleware)

		r.Get("/heartbeat", func(w http.ResponseWriter, r *http.Request) {
			remote := r.Context().Value("env").(env_manager.RemoteEnvironment)

			remote.Heartbeat()

			if remote.IsNeeded() {
				w.WriteHeader(http.StatusAccepted)
			}
			w.Write([]byte("OK"))
			return
		})

		r.Get("/ws", func(w http.ResponseWriter, r *http.Request) {
			remote := r.Context().Value("env").(env_manager.RemoteEnvironment)

			HandleWebsocketConnection(w, r, remote)
		})
	})

	return r
}
