package provider

import (
	"github.com/jcelliott/lumber"
	"github.com/nanobox-io/golang-docker-client"

	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/provider"
)

// Init initializes the docker client for the provider
func Init() error {

	// load the docker environment
	if err := provider.DockerEnv(); err != nil {
		lumber.Error("provider:Init:provider.DockerEnv(): %s", err.Error())
		return util.ErrorAppend(err, "failed to load the docker environment")
	}

	// initialize the docker client
	if err := docker.Initialize("env"); err != nil {
		lumber.Error("provider:Init:docker.Initialize()")
		return util.ErrorAppend(err, "failed to initialize the docker client")
	}

	return nil
}
