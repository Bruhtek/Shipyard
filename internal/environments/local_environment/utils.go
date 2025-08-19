package local_environment

import (
	"slices"
	"strings"
)

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

func selectIdPrefixFromList(ids []string, prefix string) string {
	matches := make([]string, 0, len(ids))
	for _, id := range ids {
		if strings.HasPrefix(id, prefix) {
			matches = append(matches, id)
		}

		if len(matches) > 1 {
			return ""
		}
	}

	if len(matches) == 1 {
		return matches[0]
	}

	return ""
}
