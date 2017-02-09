package env

import (
	"fmt"

	"github.com/jcelliott/lumber"
	"github.com/nanobox-io/golang-docker-client"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors/app"
	"github.com/nanobox-io/nanobox/processors/provider"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/locker"
	util_provider "github.com/nanobox-io/nanobox/util/provider"
)

// Destroy brings down the environment setup
func Destroy(env *models.Env) error {
	locker.LocalLock()
	defer locker.LocalUnlock()

	// init docker client
	if err := provider.Init(); err != nil {
		return util.ErrorAppend(err, "failed to init docker client")
	}

	// find apps
	apps, err := env.Apps()
	if err != nil {
		lumber.Error("env:Destroy:models.Env{ID:%s}.Apps(): %s", env.ID, err)
		return util.ErrorAppend(err, "failed to load app collection")
	}

	// destroy apps
	for _, a := range apps {

		err := app.Destroy(a)
		if err != nil {
			return util.ErrorAppend(err, "failed to remove app")
		}
	}

	// unmount the environment
	if err := Unmount(env); err != nil {
		return util.ErrorAppend(err, "failed to unmount env")
	}

	// TODO: remove folder from host /mnt/sda1/env_id
	if err := util_provider.RemoveEnvDir(env.ID); err != nil {
		// it is ok if the cleanup fails its not worth erroring here
		// return util.ErrorAppend(err, "failed to remove the environment from host")
	}

	// remove volumes
	docker.VolumeRemove(fmt.Sprintf("nanobox_%s_app", env.ID))
	docker.VolumeRemove(fmt.Sprintf("nanobox_%s_cache", env.ID))
	docker.VolumeRemove(fmt.Sprintf("nanobox_%s_mount", env.ID))
	docker.VolumeRemove(fmt.Sprintf("nanobox_%s_deploy", env.ID))
	docker.VolumeRemove(fmt.Sprintf("nanobox_%s_build", env.ID))

	// remove the environment
	if err := env.Delete(); err != nil {
		return util.ErrorAppend(err, "failed to remove env")
	}

	return nil
}
