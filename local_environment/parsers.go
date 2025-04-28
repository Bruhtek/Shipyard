package local_environment

import (
	"Shipyard/docker"
	"encoding/json"
	"log"
	"strings"
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
