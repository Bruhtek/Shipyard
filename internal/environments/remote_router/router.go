package remote_router

import (
	"Shipyard/internal/environments"
	"Shipyard/internal/environments/types"
	"context"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func RemoteMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Query().Get("key")
		if key == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if env := environments.EnvManager.GetRemoteEnv(key); env != nil {
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

	r.Route("/", func(r chi.Router) {
		r.Use(RemoteMiddleware)

		r.Get("/heartbeat", func(w http.ResponseWriter, r *http.Request) {
			remote := r.Context().Value("env").(types.RemoteEnvironment)

			remote.Heartbeat()

			if remote.IsNeeded() {
				w.WriteHeader(http.StatusAccepted)
			}
			w.Write([]byte("OK"))
			return
		})

		r.Get("/ws", func(w http.ResponseWriter, r *http.Request) {
			remote := r.Context().Value("env").(types.RemoteEnvironment)

			HandleWebsocketConnection(w, r, remote)
		})
	})

	return r
}
