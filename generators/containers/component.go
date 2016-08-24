package containers

import (
  "fmt"
  
  "github.com/nanobox-io/golang-docker-client"
  
  "github.com/nanobox-io/nanobox/models"
)

// ComponentConfig generates the container configuration for a component container
func ComponentConfig(c *models.Component, image, ip string) docker.ContainerConfig {
  config := docker.ContainerConfig{
    Name:    ComponentContainerName(c),
    Image:   image,
    Network: "virt",
    IP:      ip,
  }
  
  return config
}

// ComponentContainerName returns the name of the component container
func ComponentContainerName(c *models.Component) string {
  return fmt.Sprintf("nanobox_%s_%s", c.AppID, c.Name)
}
