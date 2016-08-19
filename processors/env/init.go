package env

import (
	"github.com/nanobox-io/golang-docker-client"

	"github.com/nanobox-io/nanobox/util/provider"
)

// Init ...
type Init struct {
}

//
func (envSetup *Init) Run() error {
	if err := provider.DockerEnv(); err != nil {
		return err
	}

	if err := docker.Initialize("env"); err != nil {
		return err
	}
	return nil
}
