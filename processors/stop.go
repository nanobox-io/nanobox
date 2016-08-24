package processors

import (
	"fmt"
	
	"github.com/jcelliott/lumber"
	
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors/app"
	"github.com/nanobox-io/nanobox/processors/env"
	"github.com/nanobox-io/nanobox/processors/provider"
)

// Stop stops the running apps, unmounts all envs, and stops the provider
func Stop() error {

	// stop all running apps
	if err := stopAllApps(); err != nil {
		return fmt.Errorf("failed to stop running apps: %s", err.Error())
	}

	// unmount envs
	if err := unmountEnvs(); err != nil {
		return fmt.Errorf("failed to unmount envs: %s", err.Error())
	}

	// stop the provider
	if err := provider.Stop(); err != nil {
		return fmt.Errorf("failed to stop the provider: %s", err.Error())
	}
	
	return nil
}

// stopAllApps stops all of the apps that are currently running
func stopAllApps() error {

	// load all the apps that think they're currently up
	apps, err := models.AllAppsByStatus("up")
	if err != nil {
		lumber.Error("stopAllApps:models.AllAppsByStatus(up): %s", err.Error())
		return fmt.Errorf("failed to load running apps: %s", err.Error())
	}

	// run the app stop on all running apps
	for _, a := range apps {
		appStop := app.Stop{App: a}
		if err := appStop.Run(); err != nil {
			return fmt.Errorf("failed to stop running app: %s", err.Error())
		}
	}
	
	return nil
}

// unmountEnvs unmounts all of the environments
func unmountEnvs() error {
	// unmount all the environments so stoping doesnt take forever
	envs, err := models.AllEnvs()
	if err != nil {
		lumber.Error("unmountEnvs:models.AllEnvs(): %s", err.Error())
		return fmt.Errorf("failed to load all envs: %s", err.Error())
	}

	for _, e := range envs {
		envUnmount := env.Unmount{Env: e}
		if err := envUnmount.Run(); err != nil {
			return fmt.Errorf("failed to unmount env: %s", err.Error())
		}
	}
	
	return nil
}
