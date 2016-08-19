package env

import (
	"fmt"

	"github.com/nanobox-io/nanobox/models"
	netfs_processor "github.com/nanobox-io/nanobox/processors/env/netfs"
	"github.com/nanobox-io/nanobox/util/provider"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/netfs"
)

// Mount ...
type Mount struct {
	Env models.Env
}

//
func (mount *Mount) Run() error {

	// mount the engine if it's a local directory
	if config.EngineDir() != "" {
		src := config.EngineDir()
		dst := fmt.Sprintf("%s%s/engine", provider.HostShareDir(), mount.Env.ID)

		// first export the env on the workstation
		if err := mount.addShare(src, dst); err != nil {
			return err
		}

		// mount the env on the provider
		if err := mount.addMount(src, dst); err != nil {
			return err
		}
	}

	// mount the app src
	src := mount.Env.Directory
	dst := fmt.Sprintf("%s%s/code", provider.HostShareDir(), mount.Env.ID)

	// first export the env on the workstation
	if err := mount.addShare(src, dst); err != nil {
		return fmt.Errorf("addShare:%s", err.Error())
	}

	// then mount the env on the provider
	if err := mount.addMount(src, dst); err != nil {
		return fmt.Errorf("addMount:%s", err.Error())
	}

	return nil
}

// addShare will add a filesystem env on the host machine
func (mount *Mount) addShare(src, dst string) error {

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

			netfsRemove := netfs_processors.Remove{src}
			if err := netfsRemove.Run(); err != nil {
				return err
			}
		}

		// now we let the provider add it's native env
		if err := provider.AddShare(src, dst); err != nil {
			return err
		}

	// check to see if native envs are currently exported. If so,
	// tear down the native env and build the netfs env
	case "netfs":
		if provider.HasShare(src, dst) {
			// native was used prior. So we need to tear it down
			if err := provider.RemoveShare(src, dst); err != nil {
				return err
			}
		}

		netfsAdd := netfs_processors.Add{src}
		if err := netfsAdd.Run(); err != nil {
			return err
		}
	}

	return nil
}

// addMount will mount a env in the nanobox guest context
func (mount *Mount) addMount(src, dst string) error {

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
			return err
		}

	// build the netfs mount
	case "netfs":
		if err := netfs.Mount(src, dst); err != nil {
			return err
		}
	}

	return nil
}
