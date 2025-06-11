package local_environment

import (
	"Shipyard/internal/docker"
	"Shipyard/internal/terminal_simple"
	"fmt"
	"github.com/rs/zerolog/log"
	"strings"
)

func (e *LocalEnvironment) GetImageCount() int {
	e.imageMutex.RLock()
	defer e.imageMutex.RUnlock()

	return len(e.images)
}

func (e *LocalEnvironment) ScanImages() {
	e.imageMutex.Lock()
	defer e.imageMutex.Unlock()

	out, err := terminal_simple.RunSimpleCommand("docker image ls --format json --no-trunc")
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
		currentImage, ok := e.images[image.ID]
		if ok && currentImage.RepoDigests != nil {
			images[num].RepoDigests = currentImage.RepoDigests
		} else {
			out, err = terminal_simple.RunSimpleCommand(
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
		return nil
	}

	return image
}
