package environments

import (
	"Shipyard/internal/docker"
	"Shipyard/internal/environments/types"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func GetStacksRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		env, ok := r.Context().Value("env").(types.LocalEnvironment)
		if !ok {
			http.Error(w, "Environment not found or not a local environment", http.StatusNotFound)
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
