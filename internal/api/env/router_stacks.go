package env

import (
	"Shipyard/internal/docker"
	"Shipyard/internal/env_manager"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func GetStacksRouter() *chi.Mux {
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

		if env.GetStackCount() == 0 {
			env.ScanStacks()
		}

		stacks := env.GetStacks()

		type StackList struct {
			Stacks map[string]*docker.Stack
			Length int
		}

		stackList := StackList{
			Stacks: stacks,
			Length: len(stacks),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(stackList); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	return r
}
