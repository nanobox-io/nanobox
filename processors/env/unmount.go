package env

import (
	"fmt"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors/env/netfs"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/display"
	"github.com/nanobox-io/nanobox/util/provider"
)

// Unmount unmounts the env shares
func Unmount(env *models.Env, keepShare bool) error {
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
		dst := fmt.Sprintf("%s%s/engine", provider.HostShareDir(), env.ID)

		// unmount the env on the provider
		if err := removeMount(src, dst); err != nil {
			display.ErrorTask()
			return fmt.Errorf("failed to remove engine mount: %s", err.Error())
		}

		// remove the share
		if !keepShare {
			if err := removeShare(src, dst); err != nil {
				display.ErrorTask()
				return fmt.Errorf("failed to remove engine share: %s", err.Error())
			}
		}
	}

	// unmount the app src
	src := env.Directory
	dst := fmt.Sprintf("%s%s/code", provider.HostShareDir(), env.ID)

	// unmount the env on the provider
	if err := removeMount(src, dst); err != nil {
		display.ErrorTask()
		return fmt.Errorf("failed to remove code mount: %s", err.Error())
	}

	// then remove the share from the workstation
	if !keepShare {
		if err := removeShare(src, dst); err != nil {
			display.ErrorTask()
			return fmt.Errorf("failed to remove code share: %s", err.Error())
		}
	}

	return nil
}

func removeMount(src, dst string) error {

	// short-circuit if the mount doesnt exist
	if !provider.HasMount(dst) {
		return nil
	}

	if err := provider.RemoveMount(src, dst); err != nil {
		return fmt.Errorf("failed to remove mount: %s", err.Error())
	}

	return nil
}

// removeShare will add a filesystem env on the host machine
func removeShare(src, dst string) error {

	// the mount type is configurable by the user
	mountType := config.Viper().GetString("mount-type")

	switch mountType {

	case "native":
		// remove the native mount using the provider
		if err := provider.RemoveShare(src, dst); err != nil {
			return fmt.Errorf("failed to remove native mount: %s", err.Error())
		}

	case "netfs":
		// remove the netfs mount using its processor
		if err := netfs.Remove(src); err != nil {
			return fmt.Errorf("failed to remove netfs share: %s", err.Error())
		}
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
		dst := fmt.Sprintf("%s%s/engine", provider.HostShareDir(), env.ID)
		
		if provider.HasMount(dst) {
			return true
		}
	}
	
	// check to see if the code is mounted
	dst := fmt.Sprintf("%s%s/code", provider.HostShareDir(), env.ID)
	return provider.HasMount(dst)
}
