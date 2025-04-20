package actions

import (
	"Shipyard/api/websocket"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func GetActionsRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		actions := websocket.ActionManager.GetActions()

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		type ActionList struct {
			Actions map[string]*websocket.Action
		}

		actionList := ActionList{
			Actions: actions,
		}

		for _, action := range actionList.Actions {
			action.Mutex.RLock()
			defer action.Mutex.RUnlock()
		}

		if err := json.NewEncoder(w).Encode(actionList); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	return r
}
