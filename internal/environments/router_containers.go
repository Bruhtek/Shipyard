package environments

import (
	"Shipyard/internal/docker"
	"Shipyard/internal/environments/types"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func GetContainersRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		env, ok := r.Context().Value("env").(types.LocalEnvironment)
		if !ok {
			http.Error(w, "Environment not found or not a local environment", http.StatusNotFound)
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
