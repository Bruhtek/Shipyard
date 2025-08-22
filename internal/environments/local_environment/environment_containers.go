package local_environment

import (
	"Shipyard/internal/docker"
	"Shipyard/internal/terminals"
	"context"
	"fmt"
	"maps"
	"slices"
	"strings"
	"time"

	"github.com/regclient/regclient"
	"github.com/regclient/regclient/types/ref"
	"github.com/rs/zerolog/log"
)

const UPDATE_CHECK_COOLDOWN = time.Hour * 6
const ERROR_UPDATE_CHECK_COOLDOWN = UPDATE_CHECK_COOLDOWN / 2
const MAX_CHECKS_PER_SCAN = 5

var RATELIMIT_LAST_WARNING = time.Time{}

const RATELIMIT_WARNING_COOLDOWN = time.Hour * 1

func (e *LocalEnvironment) ScanContainers() {
	out, err := terminals.RunSimpleCommand(ContainerLsCommand)
	if err != nil {
		log.Err(err).Msg("Error listing containers")
		return
	}

	updateCheckCount := 0

	containers := ParsePsJson([]byte(out))
	for id, container := range containers {
		currentContainer := e.GetContainer(container.ID)
		//#region Check Container ImageID
		if currentContainer != nil && currentContainer.ImageID != "" {
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

		if currentContainer != nil {
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

	e.containerMutex.Lock()
	defer e.containerMutex.Unlock()
	e.containers = make(map[string]*docker.Container)
	for _, container := range containers {
		e.containers[container.ID] = container
	}
}

func (e *LocalEnvironment) checkContainerUpdateStatus(container *docker.Container) {
	if time.Since(RATELIMIT_LAST_WARNING) < RATELIMIT_WARNING_COOLDOWN {
		// if we recently hit a rate limit, skip further checks for a while
		container.UpToDate = docker.Pending
		return
	}

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
		if strings.HasPrefix(container.Image, "sha256:") {
			// it is very possible, that the image doesn't exist,
			// or has been modified in some way during the container lifecycle
			container.UpToDate = docker.Unknown
		}

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
		// warn, since this can happen with some images - local only, for example
		log.Warn().
			Str("container-id", container.ID).
			Str("container-name", container.Name).
			Msg("Container image has no repo digests, cannot check for updates - is this a local image?")
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
			Msg("Error while checking container update status: Error getting manifest. " +
				"It's possible this is a local image, or it's behind authentication.")
		container.UpToDate = docker.Unknown
		return
	}

	digests := make([]string, 0)

	if manifest.IsList() {
		fullManifest, err := rc.ManifestGet(ctx, r)
		if err != nil {
			if strings.Contains(err.Error(), "429") {
				log.Warn().
					Msg("Hit rate limit exceeded while checking container update status - stopping further checks for 1 hour")
				RATELIMIT_LAST_WARNING = time.Now()
				container.UpToDate = docker.Pending
				return
			} else {
				log.Err(err).
					Str("container-id", container.ID).
					Str("container-name", container.Name).
					Str("image", container.Image).
					Msg("Error while checking container update status: Error getting full manifest")
				container.UpToDate = docker.Unknown
			}
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

	cont, ok := e.containers[id]
	if !ok {
		id = selectIdPrefixFromList(slices.Collect(maps.Keys(e.containers)), id)

		return e.containers[id]
	}

	return cont
}

func (e *LocalEnvironment) GetContainerCount() int {
	e.containerMutex.RLock()
	defer e.containerMutex.RUnlock()

	return len(e.containers)
}
