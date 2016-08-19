package env

import (
	"fmt"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors/app"
	"github.com/nanobox-io/nanobox/processors/provider"
	"github.com/nanobox-io/nanobox/util/locker"
)

// Destroy ...
type Destroy struct {
	Env models.Env
}

//
func (destroy *Destroy) Run() error {
	locker.LocalLock()
	defer locker.LocalUnlock()

	// we need the vm to be up and running
	providerSetup := provider.Setup{}
	if err := providerSetup.Run(); err != nil {
		return err
	}

	// find apps
	apps, err := models.AllAppsByEnv(destroy.Env.ID)
	if err != nil {
		return err
	}

	// destroy apps
	for _, a := range apps {
		appDestroy := app.Destroy{
			App: a,
		}

		err := appDestroy.Run()
		if err != nil {
			return fmt.Errorf("failed to remove app: %s", err.Error())
		}
	}

	// destroy the mounts
	return destroy.destroyMounts()
}

// destroyMounts removes the environments mounts from the provider
func (destroy *Destroy) destroyMounts() error {
	// we will not be tearing down the mounts currently
	// this is because they are required for production builds
	return nil
	// return processors.Run("app_unmount", destroy.control)
}
