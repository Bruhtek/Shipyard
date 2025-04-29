package clients

import "github.com/go-chi/chi/v5"

func GetClientsRouter() *chi.Mux {
	r := chi.NewRouter()

	return r
}
