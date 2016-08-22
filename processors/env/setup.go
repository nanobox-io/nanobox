package env

import (
	"github.com/nanobox-io/nanobox/commands/registry"
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors/app"
	"github.com/nanobox-io/nanobox/processors/provider"
	"github.com/nanobox-io/nanobox/util/config"
)

type Setup struct {
	Env models.Env
}

// Run sets up the provider and the env mounts
func (s *Setup) Run() error {

	// ensure the env data has been generated
	if err := s.Env.Generate(); err != nil {
		lumber.Error("env:Setup:Run:models:Env:Generate(): %s", err.Error())
		return fmt.Errorf("failed to initialize the env data: %s", err.Error())
	}

	

	if err := setup.setupProvider(); err != nil {
		return err
	}

	if err := setup.setupMounts(); err != nil {
		return err
	}

	// TODO: we shouldn't be doing this here
	// // if there is an environment then we should set up app
	// // if not (in the case of a build) no app setup is necessary
	// if registry.GetString("appname") != "" {
	// 	if err := setup.setupApp(); err != nil {
	// 		return err
	// 	}
	// }

	return nil
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

// // setupApp sets up the app plaftorm and data services
// func (setup *Setup) setupApp() error {
// 
// 	// setup the app
// 	appSetup := app.Setup{Env: setup.Env}
// 
// 	return appSetup.Run()
// }
