package docker

import (
	"encoding/json"
	"slices"
	"time"
)

type Network struct {
	CreatedAt time.Time
	Driver    string
	ID        string
	Internal  bool
	IPv6      bool
	Labels    map[string]string
	Name      string
	Scope     string

	Containers []*Container
}

func (n *Network) toJSON() ([]byte, error) {
	return json.Marshal(n)
}

func (n *Network) UpdateNetworkContainers(containers map[string]*Container) {
	name := n.Name
	n.Containers = make([]*Container, 0)
	for _, container := range containers {
		if slices.Contains(container.Networks, name) {
			n.Containers = append(n.Containers, container)
		}
	}
}
