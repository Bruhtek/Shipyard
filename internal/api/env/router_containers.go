package env

import (
	"Shipyard/internal/docker"
	"Shipyard/internal/env_manager"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func GetContainersRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		envI := r.Context().Value("env").(env_manager.EnvInterface)

		env, ok := envI.(env_manager.LocalEnvironment)
		if !ok {
			remote, ok := envI.(env_manager.RemoteEnvironment)
			if !ok {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			res, err := remote.GetResponse(r.URL.Path)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Error retrieving response from remote"))
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(res.Code)
			w.Write([]byte(res.Body))
			return
		}

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

	return r
}
