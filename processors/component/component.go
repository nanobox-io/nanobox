// Package service ...
package component

import (
	"fmt"

	"github.com/nanobox-io/golang-docker-client"
	"github.com/nanobox-io/nanobox-boxfile"

	"github.com/nanobox-io/nanobox/models"
)

// isComponentRunning returns true if a service is already running
func isComponentRunning(containerID string) bool {
	container, err := docker.GetContainer(containerID)

	// if the container doesn't exist then just return false
	return err == nil && container.State.Status == "running"
}

// componentImage returns the image for the component
func componentImage(component *models.Component) (string, error) {
	// fetch the env
	env, err := models.FindEnvByID(component.EnvID)
	if err != nil {
		return "", fmt.Errorf("failed to load env model: %s", err.Error())
	}

	box := boxfile.New([]byte(env.BuiltBoxfile))
	image := box.Node(component.Name).StringValue("image")

	// the only way image can be empty is if it's a platform service
	if image == "" {
		image = fmt.Sprintf("nanobox/%s", component.Name)
	}

	return image, nil
}
