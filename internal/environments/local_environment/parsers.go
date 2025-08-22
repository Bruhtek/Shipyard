package local_environment

import (
	"Shipyard/internal/docker"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

func parseCreatedAt(createdAtStr string) (time.Time, error) {
	createdAt, err := time.Parse("2006-01-02 15:04:05 -0700 MST", createdAtStr)
	if err != nil {
		createdAt, err = time.Parse("2006-01-02 15:04:05 -0700 -0700", createdAtStr)
		if err != nil {
			log.Err(err).Str("createdAt", createdAtStr).Msg("Error parsing createdAt time")
			return time.Time{}, err
		}
	}

	return createdAt, nil
}

func parseLabels(labelsStr string) (map[string]string, error) {
	labels := make(map[string]string)
	if labelsStr != "" {
		splitLabels := strings.Split(labelsStr, ",")
		prevLabelKey := ""

		for _, label := range splitLabels {
			labelSplit := strings.Split(label, "=")
			if len(labelSplit) == 2 {
				labels[labelSplit[0]] = labelSplit[1]
				prevLabelKey = labelSplit[0]
			} else if len(labelSplit) == 1 {
				previousValue, exists := labels[prevLabelKey]
				if exists {
					// If there is no = sign, we assume it's a continuation of the previous label's value
					// This is a workaround for labels that might not be formatted correctly
					labels[prevLabelKey] = previousValue + "," + labelSplit[0]
				} else {
					log.Error().Str("label", label).Msg("Invalid label format without key")
					return make(map[string]string), errors.New("invalid label format without key")
				}
			} else {
				log.Error().Str("label", label).Msg("Invalid label format")
				return make(map[string]string), errors.New("invalid label format")
			}
		}
	} else {
		return labels, nil
	}

	return labels, nil
}

const ContainerLsCommand = "docker container ls -a --no-trunc --format {{.ID}}\t{{.Image}}\t{{.Labels}}\t{{.Names}}\t{{.Networks}}\t{{.Ports}}\t{{.State}}\t{{.Status}}\t{{.CreatedAt}}\t{{.Command}}"

func ParsePsJson(jsonData []byte) []*docker.Container {
	splitData := strings.Split(string(jsonData), "\n")
	containers := make([]*docker.Container, 0)

	for _, line := range splitData {
		if line == "" {
			continue
		}

		splitLine := strings.Split(line, "\t")
		if len(splitLine) < 10 {
			log.Error().Str("line", line).Msg("Invalid container data")
			continue
		}
		createdAt, err := parseCreatedAt(splitLine[8])
		if err != nil {
			continue
		}

		labels, err := parseLabels(splitLine[2])
		if err != nil {
			log.Err(err).Str("labels", splitLine[2]).Msg("Error parsing labels")
			continue
		}

		names := strings.Split(splitLine[3], ",")
		name := names[0]

		networks := strings.Split(splitLine[4], ",")
		ports := strings.Split(splitLine[5], ", ")

		container := &docker.Container{
			ID:              splitLine[0],
			Image:           splitLine[1],
			Labels:          labels,
			Name:            name,
			Names:           names,
			Ports:           ports,
			Networks:        networks,
			State:           splitLine[6],
			Status:          splitLine[7],
			CreatedAt:       createdAt,
			Command:         strings.Trim(splitLine[9], "\""),
			ImageID:         "",
			UpToDate:        0,
			LastUpdateCheck: time.Time{},
		}

		containers = append(containers, container)
	}

	return containers
}

func ParseImageLsJson(jsonData []byte) []docker.Image {
	splitData := strings.Split(string(jsonData), "\n")
	images := make([]docker.Image, 0)

	for _, line := range splitData {
		if line == "" {
			continue
		}

		tempImage := docker.TempImage{}
		err := json.Unmarshal([]byte(line), &tempImage)
		if err != nil {
			log.Err(err).Msg("Error parsing image from JSON")
			continue
		}

		image, err := tempImage.ToImage()
		if err != nil {
			log.Err(err).Msg("Error converting temp image to image")
			continue
		}

		images = append(images, image)
	}

	return images
}

const NetworkLsCommand = "docker network ls --no-trunc --format {{.ID}};{{.Name}};{{.CreatedAt}};{{.Driver}};{{.Internal}};{{.IPv6}};{{.Scope}};{{.Labels}}"

func ParseNetworkLsJson(jsonData *string) []docker.Network {
	networksStrings := strings.Split(*jsonData, "\n")
	networks := make([]docker.Network, 0)

	for _, line := range networksStrings {
		if line == "" {
			continue
		}

		splitLine := strings.Split(line, ";")
		if len(splitLine) < 8 {
			log.Error().Str("line", line).Msg("Invalid network data")
			continue
		}
		createdAt, err := parseCreatedAt(splitLine[2])
		if err != nil {
			continue
		}

		labels, err := parseLabels(splitLine[7])

		network := docker.Network{
			ID:        splitLine[0],
			Name:      splitLine[1],
			CreatedAt: createdAt,
			Driver:    splitLine[3],
			Internal:  splitLine[4] == "true",
			IPv6:      splitLine[5] == "true",
			Scope:     splitLine[6],
			Labels:    labels,
		}

		networks = append(networks, network)
	}

	return networks
}

func ParseComposeJson(jsonData *string) []*docker.Stack {
	stacks := make([]*docker.Stack, 0)
	err := json.Unmarshal([]byte(*jsonData), &stacks)
	if err != nil {
		log.Err(err).Msg("Error parsing compose json")
		return nil
	}

	return stacks
}
