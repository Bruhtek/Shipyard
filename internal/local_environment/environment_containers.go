package local_environment

import (
	"Shipyard/internal/docker"
	"Shipyard/internal/terminals"
	"context"
	"fmt"
	"github.com/regclient/regclient"
	"github.com/regclient/regclient/types/ref"
	"github.com/rs/zerolog/log"
	"strings"
	"time"
)

const UPDATE_CHECK_COOLDOWN = time.Hour * 2
const ERROR_UPDATE_CHECK_COOLDOWN = UPDATE_CHECK_COOLDOWN / 2

func (e *LocalEnvironment) ScanContainers() {
	e.containerMutex.Lock()
	defer e.containerMutex.Unlock()

	out, err := terminals.RunSimpleCommand("docker ps -a --format json --no-trunc")
	if err != nil {
		log.Err(err).Msg("Error listing containers")
		return
	}

	containers := ParsePsJson([]byte(out))
	for id, container := range containers {
		currentContainer, ok := e.containers[container.ID]
		//#region Check Container ImageID
		if ok {
			// image ID is immutable, so we can skip the relatively expensive inspect command if we already have it
			containers[id].ImageID = currentContainer.ImageID
		} else {
			out, err = terminals.RunSimpleCommand(
				fmt.Sprintf("docker container inspect --format '{{.Image}}' %s", container.ID))
			if err != nil {
				log.Err(err).
					Str("container-id", container.ID).
					Str("container-name", container.Name).
					Msg("Error inspecting container")
				continue
			}
			containers[id].ImageID = strings.Trim(strings.TrimSpace(out), "'")
		}
		//#endregion

		if ok {
			shouldUpdate := true
			if currentContainer.UpToDate == docker.UpdateAvailable {
				shouldUpdate = false
			}
			if currentContainer.UpToDate != docker.Unknown &&
				time.Since(currentContainer.LastUpdateCheck) < UPDATE_CHECK_COOLDOWN {
				shouldUpdate = false
			}
			if currentContainer.UpToDate == docker.Error &&
				time.Since(currentContainer.LastUpdateCheck) >= ERROR_UPDATE_CHECK_COOLDOWN {
				shouldUpdate = true
			}

			if shouldUpdate {
				e.checkContainerUpdateStatus(containers[id])
			} else {
				containers[id].LastUpdateCheck = currentContainer.LastUpdateCheck
				containers[id].UpToDate = currentContainer.UpToDate
			}
		} else {
			e.checkContainerUpdateStatus(containers[id])
		}
	}

	e.containers = make(map[string]*docker.Container)
	for _, container := range containers {
		e.containers[container.ID] = container
	}
}

func (e *LocalEnvironment) checkContainerUpdateStatus(container *docker.Container) {
	container.LastUpdateCheck = time.Now()
	rc := regclient.New()

	imageRef, err := ref.New(container.Image)
	if err != nil {
		log.Err(err).Msg("Error while checking container update status: Error parsing image ref")
		container.UpToDate = docker.Error
		return
	}

	ctx := context.Background()
	manifest, err := rc.ManifestHead(ctx, imageRef)
	if err != nil {
		log.Err(err).Msg("Error while checking container update status: Error getting manifest")
		container.UpToDate = docker.Error
		return
	}

	defer rc.Close(ctx, imageRef)

	manifestDigest := manifest.GetDescriptor().Digest.String()

	log.Debug().
		Str("manifest-digest", manifestDigest).
		Str("image-id", container.ImageID).
		Msg("Checking container update status")

	if manifestDigest == container.ImageID {
		container.UpToDate = docker.UpToDate
	} else {
		container.UpToDate = docker.UpdateAvailable
	}
}

func (e *LocalEnvironment) GetContainers() map[string]*docker.Container {
	e.containerMutex.RLock()
	defer e.containerMutex.RUnlock()

	return e.containers
}

func (e *LocalEnvironment) GetContainer(id string) *docker.Container {
	e.containerMutex.RLock()
	defer e.containerMutex.RUnlock()

	return e.containers[id]
}

func (e *LocalEnvironment) GetContainerCount() int {
	e.containerMutex.RLock()
	defer e.containerMutex.RUnlock()

	return len(e.containers)
}
