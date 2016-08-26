package env

import (
	"fmt"

	"github.com/jcelliott/lumber"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors/app"
	"github.com/nanobox-io/nanobox/util/locker"
)

// Destroy brings down the environment setup
func Destroy(env *models.Env) error {
	locker.LocalLock()
	defer locker.LocalUnlock()

	// find apps
	apps, err := env.Apps()
	if err != nil {
		lumber.Error("env:Destroy:models.Env{ID:%s}.Apps(): %s", env.ID, err.Error())
		return fmt.Errorf("failed to load app collection: %s", err.Error())
	}

	// destroy apps
	for _, a := range apps {

		err := app.Destroy(a)
		if err != nil {
			return fmt.Errorf("failed to remove app: %s", err.Error())
		}
	}

	return nil
}
