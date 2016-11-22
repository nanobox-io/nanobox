package containers

import (
	"fmt"

	"github.com/nanobox-io/golang-docker-client"

	"github.com/nanobox-io/nanobox/models"
)

// ComponentConfig generates the container configuration for a component container
func ComponentConfig(componentModel *models.Component) docker.ContainerConfig {
	config := docker.ContainerConfig{
		Name:    ComponentName(componentModel),
		Image:   componentModel.Image,
		Network: "virt",
		IP:      componentModel.InternalIP,
		RestartPolicy: "no",
	}

	return config
}

// ComponentName returns the name of the component container
func ComponentName(componentModel *models.Component) string {
	return fmt.Sprintf("nanobox_%s_%s", componentModel.AppID, componentModel.Name)
}
