package local_environment

import (
	"Shipyard/internal/docker"
	"encoding/json"
	"log"
	"strings"
	"time"
)

func ParsePsJson(jsonData []byte) []docker.Container {
	splitData := strings.Split(string(jsonData), "\n")
	containers := make([]docker.Container, 0)

	for _, line := range splitData {
		if line == "" {
			continue
		}

		tempContainer := docker.TempContainer{}
		err := json.Unmarshal([]byte(line), &tempContainer)
		if err != nil {
			continue
		}

		container, err := tempContainer.ToContainer()
		if err != nil {
			log.Printf("Error parsing container: %v", err)
			continue
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
			log.Printf("Error parsing image: %v", err)
			continue
		}

		image, err := tempImage.ToImage()
		if err != nil {
			log.Printf("Error converting to image: %v", err)
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
			log.Printf("Invalid network data: %s", line)
			continue
		}
		createdAt, err := time.Parse("2006-01-02 15:04:05 -0700 MST", splitLine[2])
		if err != nil {
			log.Printf("Error parsing createdAt: %v", err)
			continue
		}

		labels := make(map[string]string)
		if splitLine[7] != "" {
			splitLabels := strings.Split(splitLine[7], ",")
			for _, label := range splitLabels {
				labelSplit := strings.Split(label, "=")
				if len(labelSplit) == 2 {
					labels[labelSplit[0]] = labelSplit[1]
				} else {
					log.Printf("Invalid label format: %s", label)
				}
			}
		}

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
