package env

import (
	"fmt"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/provider"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/data"
	"github.com/nanobox-io/nanobox/util/locker"
)

// processTeardown ...
type processTeardown struct {
	control processor.ProcessControl
}

//
func init() {
	processor.Register("env_teardown", teardownFn)
}

//
func teardownFn(control processor.ProcessControl) (processor.Processor, error) {
	return &processTeardown{control}, nil
}

//
func (teardown processTeardown) Results() processor.ProcessControl {
	return teardown.control
}

//
func (teardown *processTeardown) Process() error {

	// dont shut anything down if we are supposed to background
	if processor.DefaultControl.Debug {
		fmt.Println("leaving running because your running in debug mode")
		return nil
	}

	// if im given a environment to teardown i will tear the app down
	// in the env. If not (in the case of a build) no app is teardownable.
	if teardown.control.Env != "" {
		if err := teardown.teardownApp(); err != nil {
			return err
		}
	}

	// teardown the mounts
	return teardown.teardownMounts()
}

// teardownApp tears down the app when it's not being used
func (teardown *processTeardown) teardownApp() error {

	// establish a local app lock to ensure we're the only ones bringing down the
	// app platform. Also ensure that we release it even if we error
	locker.LocalLock()
	defer locker.LocalUnlock()

	// stop all data services
	if err := processor.Run("app_teardown", teardown.control); err != nil {
		return err
	}

	return nil
}

// teardownMounts removes the environments mounts from the provider
func (teardown *processTeardown) teardownMounts() error {

	// break early if there is still an environemnt using
	// the mounts
	if teardown.mountsInUse() {
		return nil
	}

	// establish a local app lock to ensure we're the only ones bringing
	// down the app platform. Also ensure that we release it even if we error
	locker.LocalLock()
	defer locker.LocalUnlock()

	// unmount the engine if it's a local directory
	if config.EngineDir() != "" {

		//
		src := config.EngineDir()
		dst := fmt.Sprintf("%s%s/engine", provider.HostShareDir(), config.AppName())

		// unmount the env on the provider
		if err := teardown.removeMount(src, dst); err != nil {
			return err
		}

		if err := teardown.removeShare(src, dst); err != nil {
			return err
		}
	}

	// unmount the app src
	src := config.LocalDir()
	dst := fmt.Sprintf("%s%s/code", provider.HostShareDir(), config.AppName())

	// unmount the env on the provider
	if err := teardown.removeMount(src, dst); err != nil {
		return err
	}

	// remove the share from the provider
	if err := teardown.removeShare(src, dst); err != nil {
		return err
	}

	return nil
}

func (teardown *processTeardown) mountsInUse() bool {
	app := models.App{}
	devErr := data.Get("apps", config.AppName()+"_dev", &app)
	simErr := data.Get("apps", config.AppName()+"_sim", &app)
	return devErr == nil || simErr == nil
}

// addMount will mount a env in the nanobox guest context
func (teardown *processTeardown) removeMount(src, dst string) error {

	// short-circuit if the mount doesnt exist
	if !provider.HasMount(dst) {
		return nil
	}

	return provider.RemoveMount(src, dst)
}

// addShare will add a filesystem env on the host machine
func (teardown *processTeardown) removeShare(src, dst string) error {

	// the mount type is configurable by the user
	mountType := config.Viper().GetString("mount-type")

	switch mountType {
	case "native":
		// remove the providers native share
		if err := provider.RemoveShare(src, dst); err != nil {
			return err
		}

	case "netfs":
		control := processor.ProcessControl{
			Env:          teardown.control.Env,
			Verbose:      teardown.control.Verbose,
			DisplayLevel: teardown.control.DisplayLevel,
			Meta: map[string]string{
				"path": src,
			},
		}

		if err := processor.Run("env_netfs_remove", control); err != nil {
			return err
		}
	}

	return nil
}
