package docker

import (
	"encoding/json"
	"time"
)

type Container struct {
	ID        string
	Image     string
	Labels    map[string]string
	Name      string
	Names     []string
	Ports     []string
	Networks  []string
	State     string
	Status    string
	CreatedAt time.Time
	Command   string

	// from more advanced processing
	ImageID string

	// for checking up-to-date status using regctl
	UpToDate        UpToDateStatus
	LastUpdateCheck time.Time
}

type UpToDateStatus int

const (
	Pending         UpToDateStatus = 0
	UpToDate        UpToDateStatus = 1
	UpdateAvailable UpToDateStatus = 2
	Unknown         UpToDateStatus = 3
	Pinned          UpToDateStatus = 4
)

func (c *Container) toJSON() ([]byte, error) {
	return json.Marshal(c)
}

type TempContainer struct {
	ID        string
	Image     string
	Labels    string
	Names     string
	Networks  string
	Ports     string
	State     string
	Status    string
	CreatedAt string
	Command   string
}

func AreContainersUpToDate(containers []*Container) UpToDateStatus {
	var hasUpdateAvailable = false
	var hasUnknown = false
	var allPinned = true
	var allUpToDate = true

	for _, container := range containers {
		if container.UpToDate == UpdateAvailable {
			hasUpdateAvailable = true
		}
		if container.UpToDate == Unknown {
			hasUnknown = true
		}
		if container.UpToDate != Pinned {
			allPinned = false
		}
		if container.UpToDate != UpToDate {
			allUpToDate = false
		}
	}

	if hasUnknown {
		return Unknown
	}
	if hasUpdateAvailable {
		return UpdateAvailable
	}
	if allPinned {
		return Pinned
	}
	if allUpToDate {
		return UpToDate
	}

	return Pending
}
