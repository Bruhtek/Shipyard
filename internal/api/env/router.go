package env

import (
	"Shipyard/internal/env_manager"
	"Shipyard/internal/utils"
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

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		environments := env_manager.EnvManager.GetEnvs()
		envDescriptions := make([]utils.EnvDescription, 0, len(environments))

		for _, env := range environments {
			envDescriptions = append(envDescriptions, env.GetEnvDescription())
		}

		type EnvList struct {
			Environments []utils.EnvDescription
		}
		envList := EnvList{
			Environments: envDescriptions,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(envList); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	containerRouter := GetContainersRouter()
	imagesRouter := GetImagesRouter()

	r.Route("/{environment}", func(r chi.Router) {
		r.Use(EnvironmentMiddleware)

		r.Mount("/containers", containerRouter)
		r.Mount("/images", imagesRouter)
	})

	return r
}
