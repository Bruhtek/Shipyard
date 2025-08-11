package environments

import (
	"Shipyard/internal/environments/types"
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
	"net/url"
	"sort"
	"strings"
)

func EnvironmentMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		envName := chi.URLParam(r, "environment")
		// decode the environment name to handle any URL encoding
		envName, err := url.QueryUnescape(envName)
		if err != nil {
			http.Error(w, "Invalid environment name", http.StatusBadRequest)
			return
		}

		if env := EnvManager.GetEnv(envName); env != nil {
			if env.GetEnvType() == types.EnvTypeRemote {
				HandleRemoteEnvironmentsMiddleware(next, env.(types.RemoteEnvironment)).ServeHTTP(w, r)
				return
			}

			ctx := context.WithValue(r.Context(), "env", env)

			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			http.Error(w, "Environment not found", http.StatusNotFound)
			return
		}
	})
}

func HandleRemoteEnvironmentsMiddleware(next http.Handler, env types.RemoteEnvironment) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		res, err := env.GetResponse(r.URL.Path)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error retrieving response from remote"))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(res.Code)
		w.Write([]byte(res.Body))
		return
	})
}

func GetEnvRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		environments := EnvManager.GetEnvs()
		envDescriptions := make([]types.EnvDescription, 0, len(environments))

		for _, env := range environments {
			envDescriptions = append(envDescriptions, env.GetEnvDescription())
		}

		sort.Slice(envDescriptions, func(i, j int) bool {
			return strings.ToLower(envDescriptions[i].Name) < strings.ToLower(envDescriptions[j].Name)
		})

		type EnvList struct {
			Environments []types.EnvDescription
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
	networkRouter := GetNetworksRouter()
	stackRouter := GetStacksRouter()

	r.Route("/{environment}", func(r chi.Router) {
		r.Use(EnvironmentMiddleware)

		r.Mount("/stacks", stackRouter)
		r.Mount("/containers", containerRouter)
		r.Mount("/images", imagesRouter)
		r.Mount("/networks", networkRouter)
	})

	return r
}
