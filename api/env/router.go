package env

import (
	"Shipyard/docker"
	"Shipyard/env_manager"
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func EnvironmentMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		envName := chi.URLParam(r, "environment")

		if env := env_manager.EnvManager.GetEnv(envName); env != nil {
			ctx := context.WithValue(r.Context(), "env", env)

			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			http.Error(w, "Environment not found", http.StatusNotFound)
			return
		}
	})
}

func GetEnvRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Route("/{environment}", func(r chi.Router) {
		r.Use(EnvironmentMiddleware)

		r.Route("/containers", func(r chi.Router) {
			r.Get("/", func(w http.ResponseWriter, r *http.Request) {
				env := r.Context().Value("env").(env_manager.EnvInterface)

				// if we already have containers, do not scan them again.
				// the intervals package takes care of that
				// this hugely improves the performance - 2 orders of magnitude faster
				if env.GetContainerCount() == 0 {
					env.ScanContainers()
				}

				containers := env.GetContainers()

				type ContainerList struct {
					Containers map[string]*docker.Container
					Length     int
				}

				containerList := ContainerList{
					Containers: containers,
					Length:     len(containers),
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				if err := json.NewEncoder(w).Encode(containerList); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			})
		})
	})

	return r
}
