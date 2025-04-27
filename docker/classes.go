package docker

import (
	"encoding/json"
	"strings"
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

	t, err := time.Parse("2006-01-02 15:04:05 -0700 MST", c.CreatedAt)
	if err != nil {
		return Container{}, err
	}

	container = Container{
		ID:        c.ID,
		Image:     c.Image,
		Labels:    make(map[string]string),
		Name:      strings.Split(c.Names, ",")[0],
		Names:     strings.Split(c.Names, ","),
		Ports:     strings.Split(c.Ports, ","),
		Networks:  strings.Split(c.Networks, ","),
		State:     c.State,
		Status:    c.Status,
		CreatedAt: t,
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
