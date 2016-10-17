package env

import (
	"fmt"

	"github.com/jcelliott/lumber"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors/app"
	"github.com/nanobox-io/nanobox/processors/provider"
	"github.com/nanobox-io/nanobox/util/locker"
	util_provider "github.com/nanobox-io/nanobox/util/provider"
)

// Destroy brings down the environment setup
func Destroy(env *models.Env) error {
	locker.LocalLock()
	defer locker.LocalUnlock()

	// init docker client
	if err := provider.Init(); err != nil {
		return fmt.Errorf("failed to init docker client: %s", err)
	}

	// find apps
	apps, err := env.Apps()
	if err != nil {
		lumber.Error("env:Destroy:models.Env{ID:%s}.Apps(): %s", env.ID, err)
		return fmt.Errorf("failed to load app collection: %s", err)
	}

	// destroy apps
	for _, a := range apps {

		err := app.Destroy(a)
		if err != nil {
			return fmt.Errorf("failed to remove app: %s", err)
		}
	}

	// unmount the environment
	if err := Unmount(env); err != nil {
		return fmt.Errorf("failed to unmount env: %s", err)
	}

	// TODO: remove folder from host /mnt/sda1/env_id
	if err := util_provider.RemoveEnvDir(env.ID); err != nil {
		return fmt.Errorf("failed to remove the environment from host: %s", err)
	}

	// remove the environment
	if err := env.Delete(); err != nil {
		return fmt.Errorf("failed to remove env: %s", err)
	}

	return nil
}
