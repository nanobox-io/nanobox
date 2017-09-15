package provider

import (
	"time"

	"github.com/jcelliott/lumber"
	"github.com/nanobox-io/golang-docker-client"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/provider"
)

// Init initializes the docker client for the provider
func Init() error {
	// load the docker environment
	if err := provider.DockerEnv(); err != nil {
		lumber.Error("provider:Init:provider.DockerEnv(): %s", err.Error())
		return util.ErrorAppend(util.ErrorQuiet(err), "failed to load the docker environment")
	}

	// initialize the docker client
	if err := docker.Initialize("env"); err != nil {
		lumber.Error("provider:Init:docker.Initialize()")
		return util.ErrorAppend(util.ErrorQuiet(err), "failed to initialize the docker client")
	}

	// make sure we have the default ip
	providerModel, _ := models.LoadProvider()
	if err := setDefaultIP(providerModel); err != nil {
		return util.ErrorAppend(err, "failed to setup the provider network")
	}

	checkFunc := func() error {
		_, err := docker.ContainerList()
		return err
	}

	// confirm it is up and working
	if err := util.Retry(checkFunc, 20, time.Second); err != nil {
		return util.Errorf("unable to communicate with Docker")
	}

	return nil
}
