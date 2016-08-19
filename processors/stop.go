package processors

import (
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors/app"
	"github.com/nanobox-io/nanobox/processors/env"
	"github.com/nanobox-io/nanobox/processors/provider"
)

// Stop ...
type Stop struct {
}

//
func (stop Stop) Run() error {

	// stop all running environments
	if err := stop.stopAllApps(); err != nil {
		return err
	}

	if err := stop.unmountEnvs(); err != nil {
		return err
	}

	// run a provider setup
	providerStop := provider.Stop{}
	return providerStop.Run()
}

// stop all of the apps that are currently up
func (stop Stop) stopAllApps() error {

	// run the app stop on all running apps
	for _, a := range upApps() {

		appStop := app.Stop{
			App: a,
		}

		if err := appStop.Run(); err != nil {
			return err
		}

	}
	return nil
}

func (stop Stop) unmountEnvs() error {
	// unmount all the environments so stoping doesnt take forever
	envs, err := models.AllEnvs()
	if err != nil	{
		return err
	}

	for _, e := range envs {

		envUnmount := env.Unmount{
			Env: e,
		}

		if err := envUnmount.Run(); err != nil {
			return err
		}

	}
	return nil	
}
