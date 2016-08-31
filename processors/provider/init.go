package provider

import (
	"fmt"

	"github.com/jcelliott/lumber"
	"github.com/nanobox-io/golang-docker-client"

	"github.com/nanobox-io/nanobox/util/provider"
)

// Init initializes the docker client for the provider
func Init() error {
	
	// load the docker environment
	if err := provider.DockerEnv(); err != nil {
		lumber.Error("provider:Init:provider.DockerEnv(): %s", err.Error())
		return fmt.Errorf("failed to load the docker environment: %s", err.Error())
	}

	// initialize the docker client
	if err := docker.Initialize("env"); err != nil {
		lumber.Error("provider:Init:docker.Initialize(): %s", err.Error())
		return fmt.Errorf("failed to initialize the docker client: %s", err.Error())
	}

	return nil
}
