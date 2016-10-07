package env

import (
	"fmt"
	"path/filepath"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/display"
	"github.com/nanobox-io/nanobox/util/provider"
)

// Unmount unmounts the env shares
func Unmount(env *models.Env) error {
	// break early if we're not mounted
	if !mounted(env) {
		return nil
	}

	// break early if there is still an environemnt using the mounts
	if mountsInUse(env) {
		return nil
	}

	display.StartTask(env.Name)
	defer display.StopTask()

	// unmount the engine if it's a local directory
	if config.EngineDir() != "" {
		src := config.EngineDir()
		dst := filepath.Join(provider.HostShareDir(), env.ID, "engine")

		// unmount the env on the provider
		if err := provider.RemoveMount(src, dst); err != nil {
			display.ErrorTask()
			return fmt.Errorf("failed to remove engine mount: %s", err.Error())
		}

	}

	// unmount the app src
	src := env.Directory
	dst := filepath.Join(provider.HostShareDir(), env.ID, "code")

	// unmount the env on the provider
	if err := provider.RemoveMount(src, dst); err != nil {
		display.ErrorTask()
		return fmt.Errorf("failed to remove code mount: %s", err.Error())
	}

	return nil
}


// mountsInUse returns true if any of the env's apps are running
func mountsInUse(env *models.Env) bool {
	devApp, _ := models.FindAppBySlug(env.ID, "dev")
	simApp, _ := models.FindAppBySlug(env.ID, "sim")
	return devApp.Status == "up" || simApp.Status == "up"
}

// returns true if the app or engine is mounted
func mounted(env *models.Env) bool {

	// if the engine is mounted, check that
	if config.EngineDir() != "" {
		dst := filepath.Join(provider.HostShareDir(), env.ID, "engine")

		if provider.HasMount(dst) {
			return true
		}
	}

	// check to see if the code is mounted
	dst := filepath.Join(provider.HostShareDir(), env.ID, "code")
	return provider.HasMount(dst)
}
