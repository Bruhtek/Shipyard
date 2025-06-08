package env

import (
	"Shipyard/internal/docker"
	"Shipyard/internal/env_manager"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func GetNetworksRouter() *chi.Mux {
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
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(res)
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
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(res)
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
