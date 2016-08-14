package env

import (
	"github.com/nanobox-io/nanobox/processor/provider"
	"github.com/nanobox-io/nanobox/processor/platform"
	"github.com/nanobox-io/nanobox/processor/app"
	"github.com/nanobox-io/nanobox/processor/component"
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/commands/registry"
)

// Setup ...
type Setup struct {
	Env models.Env
}

//
func (setup *Setup) Run() error {

	setup.setupEnv()

	if err := setup.setupProvider(); err != nil {
		return err
	}

	if err := setup.setupMounts(); err != nil {
		return err
	}

	// if there is an environment then we should set up app
	// if not (in the case of a build) no app setup is necessary
	if registry.GetString("appname") != "" {
		if err := setup.setupApp(); err != nil {
			return err
		}
	}

	return nil
}

// get the environment data
func (setup *Setup) setupEnv() error {
	setup.Env.ID = config.EnvID()
	setup.Env.Directory = config.LocalDir()
	setup.Env.Name = config.LocalDirName()
	return setup.Env.Save()
}

// setupProvider sets up the provider
func (setup *Setup) setupProvider() error {
	pSetup := provider.Setup{}

	return pSetup.Run()
}

// setupMounts will add the envs and mounts for this app
func (setup *Setup) setupMounts() error {
	mount := Mount{setup.Env}
	return mount.Run()
}

// setupApp sets up the app plaftorm and data services
func (setup *Setup) setupApp() error {

	// setup the app
	appSetup := app.Setup{Env: setup.Env}
	if err := appSetup.Run(); err != nil {
		return err
	}

	// clean up after any possible failures in a previous deploy
	componentClean := component.Clean{App: appSetup.App}
	if err := componentClean.Run(); err != nil {
		return err
	}

	// setup the platform services
	platformSetup := platform.Setup{App: appSetup.App}
	return platformSetup.Run()
}
