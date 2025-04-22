package actions

import (
	"Shipyard/api/websocket"
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
	"time"
)

func ActionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		actionId := chi.URLParam(r, "actionId")

		if action, ok := websocket.ActionManager.GetAction(actionId); ok {
			ctx := context.WithValue(r.Context(), "action", action)

			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			http.Error(w, "Action not found", http.StatusNotFound)
			return
		}
	})
}

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

	r.Route("/{actionId}", func(r chi.Router) {
		r.Use(ActionMiddleware)

		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			action := r.Context().Value("action").(*websocket.Action)

			action.Mutex.RLock()
			defer action.Mutex.RUnlock()

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(action); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		})

		r.Delete("/", func(w http.ResponseWriter, r *http.Request) {
			action := r.Context().Value("action").(*websocket.Action)

			res := action.Cancel()

			if res {
				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "application/json")

				go websocket.ActionManager.DeleteFinishedAction(action, time.Millisecond*500)

				action.Mutex.RLock()
				defer action.Mutex.RUnlock()

				if err := json.NewEncoder(w).Encode(action); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			} else {
				http.Error(w, "Failed to cancel action", http.StatusInternalServerError)
				return
			}
		})
	})

	return r
}
