// Package service ...
package component

import (
	"github.com/nanobox-io/golang-docker-client"
)

// isComponentRunning returns true if a service is already running
func isComponentRunning(containerID string) bool {
	container, err := docker.GetContainer(containerID)

	// if the container doesn't exist then just return false
	return err == nil && container.State.Status == "running"
}
