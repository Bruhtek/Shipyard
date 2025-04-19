package docker

import (
	"encoding/json"
	"strings"
)

type Container struct {
	ID        string
	Image     string
	Labels    map[string]string
	Names     []string
	Ports     []string
	Networks  []string
	State     string
	Status    string
	CreatedAt string
	Command   string
}

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

func (c *TempContainer) ToContainer() (container Container, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()

	container = Container{
		ID:        c.ID,
		Image:     c.Image,
		Labels:    make(map[string]string),
		Names:     strings.Split(c.Names, ","),
		Ports:     strings.Split(c.Ports, ","),
		Networks:  strings.Split(c.Networks, ","),
		State:     c.State,
		Status:    c.Status,
		CreatedAt: c.CreatedAt,
		Command:   c.Command,
	}

	// split labels by a comma, then split by an equal sign
	labels := strings.Split(c.Labels, ",")
	for _, label := range labels {
		labelSplit := strings.Split(label, "=")
		if len(labelSplit) != 2 {
			continue
		}
		container.Labels[labelSplit[0]] = labelSplit[1]
	}

	return container, nil
}
