package env

import (
	"fmt"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processor/env/netfs"
	"github.com/nanobox-io/nanobox/provider"
	"github.com/nanobox-io/nanobox/util/config"
)

// Unmount ...
type Unmount struct {
	Env models.Env
}

//
func (unmount *Unmount) Run() error {
	// break early if there is still an environemnt using
	// the mounts
	if unmount.mountsInUse() {
		return nil
	}

	// mount the engine if it's a local directory
	if config.EngineDir() != "" {
		src := config.EngineDir()
		dst := fmt.Sprintf("%s%s/engine", provider.HostShareDir(), unmount.Env.ID)

		// mount the env on the provider
		if err := unmount.removeMount(src, dst); err != nil {
			return err
		}

		// first export the env on the workstation
		if err := unmount.removeShare(src, dst); err != nil {
			return err
		}

	}

	// mount the app src
	src := unmount.Env.Directory
	dst := fmt.Sprintf("%s%s/code", provider.HostShareDir(), unmount.Env.ID)

	// then mount the env on the provider
	if err := unmount.removeMount(src, dst); err != nil {
		return fmt.Errorf("removeMount:%s", err.Error())
	}

	// first export the env on the workstation
	if err := unmount.removeShare(src, dst); err != nil {
		return fmt.Errorf("removeShare:%s", err.Error())
	}

	return nil
}

func (unmount *Unmount) removeMount(src, dst string) error {

	// short-circuit if the mount doesnt exist
	if !provider.HasMount(dst) {
		return nil
	}

	return provider.RemoveMount(src, dst)
}

// removeShare will add a filesystem env on the host machine
func (unmount *Unmount) removeShare(src, dst string) error {

	// the mount type is configurable by the user
	mountType := config.Viper().GetString("mount-type")

	switch mountType {

	case "native":
		// remove the native mount using the privider
		if err := provider.RemoveShare(src, dst); err != nil {
			return err
		}

	case "netfs":
		// remove the netfs mount using its processor

		netfsRemove := netfs.Remove{src}
		if err := netfsRemove.Run(); err != nil {
			return err
		}

	}

	return nil
}

func (unmount *Unmount) mountsInUse() bool {
	devApp, _ := models.FindAppBySlug(unmount.Env.ID, "dev")
	simApp, _ := models.FindAppBySlug(unmount.Env.ID, "sim")
	return devApp.Status == "up" || simApp.Status == "up"
}
