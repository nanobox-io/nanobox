package env

import (
	"fmt"

	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/provider"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/netfs"
)

// processEnvMount ...
type processEnvMount struct {
	control processor.ProcessControl
}

//
func init() {
	processor.Register("env_mount", envMountFn)
}

//
func envMountFn(control processor.ProcessControl) (processor.Processor, error) {
	// control.Meta["processEnvMount-control"]

	// do some control validation check on the meta for the flags and make sure they
	// work

	return &processEnvMount{control: control}, nil
}

//
func (envMount processEnvMount) Results() processor.ProcessControl {
	return envMount.control
}

//
func (envMount *processEnvMount) Process() error {

	// mount the engine if it's a local directory
	if config.EngineDir() != "" {
		src := config.EngineDir()
		dst := fmt.Sprintf("%s%s/engine", provider.HostShareDir(), config.AppID())

		// first export the env on the workstation
		if err := envMount.addShare(src, dst); err != nil {
			return err
		}

		// mount the env on the provider
		if err := envMount.addMount(src, dst); err != nil {
			return err
		}
	}

	// mount the app src
	src := config.LocalDir()
	dst := fmt.Sprintf("%s%s/code", provider.HostShareDir(), config.AppID())

	// first export the env on the workstation
	if err := envMount.addShare(src, dst); err != nil {
		return fmt.Errorf("addShare:%s", err.Error())
	}

	// then mount the env on the provider
	if err := envMount.addMount(src, dst); err != nil {
		return fmt.Errorf("addMount:%s", err.Error())
	}

	return nil
}

// addShare will add a filesystem env on the host machine
func (envMount *processEnvMount) addShare(src, dst string) error {

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

			control := processor.ProcessControl{
				Env:          envMount.control.Env,
				Verbose:      envMount.control.Verbose,
				DisplayLevel: envMount.control.DisplayLevel,
				Meta: map[string]string{
					"path": src,
				},
			}

			if err := processor.Run("env_netfs_remove", control); err != nil {
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

		control := processor.ProcessControl{
			Env:          envMount.control.Env,
			Verbose:      envMount.control.Verbose,
			DisplayLevel: envMount.control.DisplayLevel,
			Meta: map[string]string{
				"path": src,
			},
		}

		if err := processor.Run("env_netfs_add", control); err != nil {
			return err
		}
	}

	return nil
}

// addMount will mount a env in the nanobox guest context
func (envMount *processEnvMount) addMount(src, dst string) error {

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
