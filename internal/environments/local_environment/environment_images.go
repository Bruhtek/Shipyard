package local_environment

import (
	"Shipyard/internal/docker"
	"Shipyard/internal/terminals"
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/rs/zerolog/log"
)

func (e *LocalEnvironment) GetImageCount() int {
	e.imageMutex.RLock()
	defer e.imageMutex.RUnlock()

	return len(e.images)
}

func (e *LocalEnvironment) ScanImages() {
	out, err := terminals.RunSimpleCommand("docker image ls --format json --no-trunc")
	if err != nil {
		log.Err(err).Msg("Error listing images")
		return
	}

	// TODO: check if an image is dangling
	//_, err = terminals.RunSimpleCommand("docker images -f dangling=true -q --no-trunc")
	//if err != nil {
	//	return
	//}

	images := ParseImageLsJson([]byte(out))
	for num, image := range images {
		currentImage := e.GetImage(image.ID)
		if currentImage != nil && currentImage.RepoDigests != nil {
			// repo digests are immutable, so we can skip the relatively expensive inspect command if we already have it
			images[num].RepoDigests = currentImage.RepoDigests
		} else {
			out, err = terminals.RunSimpleCommand(
				fmt.Sprintf("docker image inspect --format {{.RepoDigests}} %s", image.ID))
			if err != nil {
				log.Err(err).
					Str("image-id", image.ID).
					Str("image-repository", image.Repository).
					Str("image-tag", image.Tag).
					Msg("Error inspecting image")
				continue
			}
			processedOut := strings.Split(strings.Trim(strings.TrimSpace(out), "[]"), ",")
			images[num].RepoDigests = make([]string, len(processedOut))
			for i, digest := range processedOut {
				images[num].RepoDigests[i] = strings.Trim(strings.TrimSpace(digest), "'\"")
			}
		}
	}

	e.imageMutex.Lock()
	defer e.imageMutex.Unlock()
	e.images = make(map[string]*docker.Image)
	for _, image := range images {
		e.images[image.ID] = &image
	}

	ids := make([]string, 0)
	for id := range e.images {
		ids = append(ids, id)
	}
	usedIds := e.getUsedImageIds(ids)
	for _, id := range usedIds {
		e.images[id].Used = true
	}

	//danglignIds = strings.TrimSpace(danglignIds)
	//danglignIdsList := strings.Split(danglignIds, "\n")
	//for _, id := range danglignIdsList {
	//	id = strings.Trim(strings.TrimSpace(id), "'")
	//	if id == "" {
	//		continue
	//	}
	//
	//}
}

func (e *LocalEnvironment) GetImages() map[string]*docker.Image {
	e.imageMutex.RLock()
	defer e.imageMutex.RUnlock()

	return e.images
}

func (e *LocalEnvironment) GetImage(id string) *docker.Image {
	e.imageMutex.RLock()
	defer e.imageMutex.RUnlock()

	image, ok := e.images[id]
	if !ok {
		if !strings.HasPrefix(id, "sha256:") {
			id = "sha256:" + id
		}
		id = selectIdPrefixFromList(slices.Collect(maps.Keys(e.images)), id)

		return e.images[id]
	}

	return image
}
