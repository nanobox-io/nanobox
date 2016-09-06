package env

import (
	"fmt"

	"github.com/nanobox-io/nanobox/models"
	netfs_processors "github.com/nanobox-io/nanobox/processors/env/netfs"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/display"
	"github.com/nanobox-io/nanobox/util/netfs"
	"github.com/nanobox-io/nanobox/util/provider"
)

// Mount sets up the env mounts
func Mount(env *models.Env) error {
	display.StartTask("Mounting codebase")
	defer display.StopTask()

	// mount the engine if it's a local directory
	if config.EngineDir() != "" {
		src := config.EngineDir()
		dst := fmt.Sprintf("%s%s/engine", provider.HostShareDir(), env.ID)

		// first export the env on the workstation
		if err := addShare(src, dst); err != nil {
			display.ErrorTask()
			return fmt.Errorf("failed to export engine share: %s", err.Error())
		}

		// mount the env on the provider
		if err := addMount(src, dst); err != nil {
			display.ErrorTask()
			return fmt.Errorf("failed to mount the engine share on the provider: %s", err.Error())
		}
	}

	// mount the app src
	src := env.Directory
	dst := fmt.Sprintf("%s%s/code", provider.HostShareDir(), env.ID)

	// first export the env on the workstation
	if err := addShare(src, dst); err != nil {
		display.ErrorTask()
		return fmt.Errorf("failed to export code share: %s", err.Error())
	}

	// then mount the env on the provider
	if err := addMount(src, dst); err != nil {
		display.ErrorTask()
		return fmt.Errorf("failed to mount the code share on the provider: %s", err.Error())
	}

	return nil
}

// addShare will add a filesystem env on the host machine
func addShare(src, dst string) error {

	// the mount type is configurable by the user
	mountType := config.Viper().GetString("mount-type")

	// todo: we should display a warning when using native about performance

	// since vm.mount is configurable, it's possible and even likely that a
	// machine may already have mounts configured. For each mount type we'll
	// need to check if an existing mount needs to be undone before continuing
	switch mountType {

	// check to see if netfs is currently configured. If it is then tear it down
	// and build the native env
	case "native":
		if netfs.Exists(src) {
			// netfs was used prior. So we need to tear it down.

			if err := netfs_processors.Remove(src); err != nil {
				return fmt.Errorf("failed to remove netfs share: %s", err.Error())
			}
		}

		// now we let the provider add it's native env
		if err := provider.AddShare(src, dst); err != nil {
			return fmt.Errorf("failed to add native share: %s", err.Error())
		}

	// check to see if native envs are currently exported. If so,
	// tear down the native env and build the netfs env
	case "netfs":
		if provider.HasShare(src, dst) {
			// native was used prior. So we need to tear it down
			if err := provider.RemoveShare(src, dst); err != nil {
				return fmt.Errorf("failed to remove native share: %s", err.Error())
			}
		}

		if err := netfs_processors.Add(src); err != nil {
			return fmt.Errorf("failed to add netfs share: %s", err.Error())
		}
	}

	return nil
}

// addMount will mount a env in the nanobox guest context
func addMount(src, dst string) error {

	// short-circuit if the mount already exists
	if provider.HasMount(dst) {
		return nil
	}

	// the mount type is configurable by the user
	mountType := config.Viper().GetString("mount-type")

	switch mountType {

	// build the native mount
	case "native":
		if err := provider.AddMount(src, dst); err != nil {
			return fmt.Errorf("failed to mount native share: %s", err.Error())
		}

	// build the netfs mount
	case "netfs":
		if err := netfs.Mount(src, dst); err != nil {
			return fmt.Errorf("failed to mount netfs share: %s", err.Error())
		}
	}

	return nil
}
