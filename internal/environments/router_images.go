package environments

import (
	"Shipyard/internal/docker"
	"Shipyard/internal/environments/types"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func GetImagesRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		env, ok := r.Context().Value("env").(types.LocalEnvironment)
		if !ok {
			http.Error(w, "Environment not found or not a local environment", http.StatusNotFound)
			return
		}

		images := env.GetImages()

		type ImageList struct {
			Images map[string]*docker.Image
			Length int
		}
		imageList := ImageList{
			Images: images,
			Length: len(images),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(imageList); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	r.Get("/{imageID}", func(w http.ResponseWriter, r *http.Request) {
		env, ok := r.Context().Value("env").(types.LocalEnvironment)
		if !ok {
			http.Error(w, "Environment not found or not a local environment", http.StatusNotFound)
			return
		}

		imageID := chi.URLParam(r, "imageID")
		image := env.GetImage(imageID)
		if image == nil {
			http.Error(w, "Image not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(image); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	return r
}
