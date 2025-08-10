package local_environment

import "slices"

func (e *LocalEnvironment) getUsedImageIds(ids []string) []string {
	usedIds := make([]string, 0)
	containers := e.GetContainers()

	containerImageIds := make([]string, 0)
	for _, container := range containers {
		containerImageIds = append(containerImageIds, container.ImageID)
	}

	for _, id := range ids {
		if slices.Contains(containerImageIds, id) {
			usedIds = append(usedIds, id)
		}
	}

	return usedIds
}
