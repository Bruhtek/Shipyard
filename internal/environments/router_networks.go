package environments

import (
	"Shipyard/internal/docker"
	"Shipyard/internal/environments/types"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func GetNetworksRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		env, ok := r.Context().Value("env").(types.LocalEnvironment)
		if !ok {
			http.Error(w, "Environment not found or not a local environment", http.StatusNotFound)
			return
		}

		if env.GetNetworkCount() == 0 {
			env.ScanNetworks()
		}

		networks := env.GetNetworks()

		type NetworkList struct {
			Networks map[string]*docker.Network
			Length   int
		}
		networkList := NetworkList{
			Networks: networks,
			Length:   len(networks),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(networkList); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	r.Get("/{networkIDorName}", func(w http.ResponseWriter, r *http.Request) {
		env, ok := r.Context().Value("env").(types.LocalEnvironment)
		if !ok {
			http.Error(w, "Environment not found or not a local environment", http.StatusNotFound)
			return
		}

		networkIDOrName := chi.URLParam(r, "networkIDorName")

		if env.GetNetworkCount() == 0 {
			env.ScanNetworks()
		}

		network := env.GetNetwork(networkIDOrName)
		if network == nil {
			http.Error(w, "Network not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(network); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	return r
}
