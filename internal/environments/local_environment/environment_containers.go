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
const MAX_CHECKS_PER_SCAN = 5

func (e *LocalEnvironment) ScanContainers() {
	e.containerMutex.Lock()
	defer e.containerMutex.Unlock()

	out, err := terminals.RunSimpleCommand(ContainerLsCommand)
	if err != nil {
		log.Err(err).Msg("Error listing containers")
		return
	}

	updateCheckCount := 0

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

		if updateCheckCount >= MAX_CHECKS_PER_SCAN {
			continue
		}

		if ok {
			shouldUpdate := true
			if currentContainer.UpToDate == docker.UpdateAvailable {
				shouldUpdate = false
			}
			if currentContainer.UpToDate != docker.Pending &&
				time.Since(currentContainer.LastUpdateCheck) < UPDATE_CHECK_COOLDOWN {
				shouldUpdate = false
			}
			if currentContainer.UpToDate == docker.Unknown &&
				time.Since(currentContainer.LastUpdateCheck) >= ERROR_UPDATE_CHECK_COOLDOWN {
				shouldUpdate = true
			}

			if shouldUpdate {
				updateCheckCount++
				e.checkContainerUpdateStatus(containers[id])
			} else {
				containers[id].LastUpdateCheck = currentContainer.LastUpdateCheck
				containers[id].UpToDate = currentContainer.UpToDate
			}
		} else {
			updateCheckCount++
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
	if strings.Contains(container.Image, "@sha256:") {
		// container image is pinned to a specific digest, no need to check for updates
		container.UpToDate = docker.Pinned
		return
	}

	imageId := container.ImageID
	if imageId == "" {
		log.Warn().
			Str("container-id", container.ID).
			Str("container-name", container.Name).
			Msg("Container has no image ID, cannot check for updates")
		container.UpToDate = docker.Unknown
		return
	}

	image := e.GetImage(imageId)
	if image == nil {
		log.Warn().
			Str("container-id", container.ID).
			Str("container-name", container.Name).
			Msg("Container image not found, cannot check for updates")
		return
	}
	if image.RepoDigests == nil {
		log.Warn().
			Str("container-id", container.ID).
			Str("container-name", container.Name).
			Msg("Container image has no repo digests, cannot check for updates")
		return
	}

	if len(image.RepoDigests) == 0 {
		// debug, not warn, since this can happen with some images - local only, for example
		log.Debug().
			Str("container-id", container.ID).
			Str("container-name", container.Name).
			Msg("Container image has no repo digests, cannot check for updates")
		container.UpToDate = docker.Unknown
		return
	}

	ctx := context.Background()
	rc := regclient.New()

	r, err := ref.New(container.Image)
	if err != nil {
		log.Err(err).
			Str("container-id", container.ID).
			Str("container-name", container.Name).
			Str("image", container.Image).
			Msg("Error while checking container update status: Error parsing image ref")
		container.UpToDate = docker.Unknown
		return
	}

	defer rc.Close(ctx, r)

	manifest, err := rc.ManifestHead(ctx, r, regclient.WithManifestRequireDigest())
	if err != nil {
		log.Err(err).
			Str("container-id", container.ID).
			Str("container-name", container.Name).
			Str("image", container.Image).
			Msg("Error while checking container update status: Error getting manifest")
		container.UpToDate = docker.Unknown
		return
	}

	digests := make([]string, 0)

	if manifest.IsList() {
		fullManifest, err := rc.ManifestGet(ctx, r)
		if err != nil {
			log.Err(err).
				Str("container-id", container.ID).
				Str("container-name", container.Name).
				Str("image", container.Image).
				Msg("Error while checking container update status: Error getting full manifest")
			container.UpToDate = docker.Unknown
			return
		}

		manifests, err := fullManifest.GetManifestList()
		if err != nil {
			log.Err(err).
				Str("container-id", container.ID).
				Str("container-name", container.Name).
				Str("image", container.Image).
				Msg("Error while checking container update status: Error getting manifest list")
			container.UpToDate = docker.Unknown
			return
		}
		for _, m := range manifests {
			digests = append(digests, m.Digest.String())
		}
	}

	digests = append(digests, manifest.GetDescriptor().Digest.String())

	log.Debug().
		Str("container-id", container.ID).
		Str("container-name", container.Name).
		Strs("digests", digests).
		Str("image", container.Image).
		Strs("repo-digests", image.RepoDigests).
		Msg("Checking container update status")

	isInDigests := false
	for _, digest := range image.RepoDigests {
		for _, mDigest := range digests {
			if strings.Contains(digest, mDigest) {
				isInDigests = true
				break
			}
		}
		if isInDigests {
			break
		}
	}
	if isInDigests {
		log.Debug().
			Str("container-id", container.ID).
			Str("container-name", container.Name).
			Str("image", container.Image).
			Msg("Container is up to date")
		container.UpToDate = docker.UpToDate
	} else {
		log.Debug().
			Str("container-id", container.ID).
			Str("container-name", container.Name).
			Str("image", container.Image).
			Msg("Container has update available")
		container.UpToDate = docker.UpdateAvailable
	}

	//imageRef, err := ref.New(container.Image)
	//if err != nil {
	//	log.Err(err).Msg("Error while checking container update status: Error parsing image ref")
	//	container.UpToDate = docker.Error
	//	return
	//}
	//
	//ctx := context.Background()
	//manifest, err := rc.ManifestHead(ctx, imageRef)
	//if err != nil {
	//	log.Err(err).Msg("Error while checking container update status: Error getting manifest")
	//	container.UpToDate = docker.Error
	//	return
	//}
	//
	//defer rc.Close(ctx, imageRef)
	//
	//manifestDigest := manifest.GetDescriptor().Digest.String()
	//
	//log.Debug().
	//	Str("manifest-digest", manifestDigest).
	//	Str("image-id", container.ImageID).
	//	Msg("Checking container update status")
	//
	//if manifestDigest == container.ImageID {
	//	container.UpToDate = docker.UpToDate
	//} else {
	//	container.UpToDate = docker.UpdateAvailable
	//}

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
