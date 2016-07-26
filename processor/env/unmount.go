package env

import (
	"fmt"
	"strings"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/provider"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/data"
)

// processEnvUnmount ...
type processEnvUnmount struct {
	control processor.ProcessControl
	app     models.App
}

//
func init() {
	processor.Register("env_unmount", envUnmountFn)
}

// returns the unmounting procedure
// Because this can be called as part of a nanobox stop
// we will need to accept a 'app_name' meta data and act accordingly:
//    the apps directory needs to come from the app model
//    the directories mounted need to be from the app root name (minus sims or devs)
func envUnmountFn(control processor.ProcessControl) (processor.Processor, error) {
	// control.Meta["processEnvUnmount-control"]

	envUnmount := &processEnvUnmount{control: control}
	return envUnmount, envUnmount.validateMeta()
}

// control validation as well as setting reasonable defaults for
// directory, app_name and app_root
func (envUnmount *processEnvUnmount) validateMeta() error {
	if envUnmount.control.Meta["directory"] == "" {
		envUnmount.control.Meta["directory"] = config.LocalDir()
	}

	// set the name of the app if we are not given one
	if envUnmount.control.Meta["app_name"] == "" {
		envUnmount.control.Meta["app_name"] = fmt.Sprintf("%s_%s", config.AppID(), envUnmount.control.Env)
	}

	// Get the root of the app (without _sim or _dev)
	if strings.HasSuffix(envUnmount.control.Meta["app_name"], "_dev") {
		envUnmount.control.Meta["app_root"] = strings.Replace(envUnmount.control.Meta["app_name"], "_dev", "", -1)
	}

	if strings.HasSuffix(envUnmount.control.Meta["app_name"], "_sim") {
		envUnmount.control.Meta["app_root"] = strings.Replace(envUnmount.control.Meta["app_name"], "_sim", "", -1)
	}

	if envUnmount.control.Meta["app_root"] == "" {
		return fmt.Errorf("i could not find a valid app root")
	}

	return nil
}

//
func (envUnmount processEnvUnmount) Results() processor.ProcessControl {
	return envUnmount.control
}

//
func (envUnmount *processEnvUnmount) Process() error {
	// break early if there is still an environemnt using
	// the mounts
	if envUnmount.mountsInUse() {
		return nil
	}

	if err := envUnmount.loadApp(); err != nil {
		return err
	}

	// mount the engine if it's a local directory
	if config.EngineDir() != "" {
		src := config.EngineDir()
		dst := fmt.Sprintf("%s%s/engine", provider.HostShareDir(), envUnmount.control.Meta["app_root"])

		// mount the env on the provider
		if err := envUnmount.removeMount(src, dst); err != nil {
			return err
		}

		// first export the env on the workstation
		if err := envUnmount.removeShare(src, dst); err != nil {
			return err
		}

	}

	// mount the app src
	src := envUnmount.control.Meta["directory"]
	dst := fmt.Sprintf("%s%s/code", provider.HostShareDir(), envUnmount.control.Meta["app_root"])

	// then mount the env on the provider
	if err := envUnmount.removeMount(src, dst); err != nil {
		return fmt.Errorf("removeMount:%s", err.Error())
	}

	// first export the env on the workstation
	if err := envUnmount.removeShare(src, dst); err != nil {
		return fmt.Errorf("removeShare:%s", err.Error())
	}

	return nil
}

func (envUnmount *processEnvUnmount) removeMount(src, dst string) error {

	// short-circuit if the mount doesnt exist
	if !provider.HasMount(dst) {
		return nil
	}

	return provider.RemoveMount(src, dst)
}

// removeShare will add a filesystem env on the host machine
func (envUnmount *processEnvUnmount) removeShare(src, dst string) error {

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

		control := processor.ProcessControl{
			Env:          envUnmount.control.Env,
			Verbose:      envUnmount.control.Verbose,
			DisplayLevel: envUnmount.control.DisplayLevel,
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

// loadApp loads the app from the db
func (envUnmount *processEnvUnmount) loadApp() error {
	// the app might not exist yet, so let's not return the error if it fails
	return data.Get("apps", envUnmount.control.Meta["name"], &envUnmount.app)
}

func (envUnmount *processEnvUnmount) mountsInUse() bool {
	devApp := models.App{}
	simApp := models.App{}
	data.Get("apps", envUnmount.control.Meta["app_root"]+"_dev", &devApp)
	data.Get("apps", envUnmount.control.Meta["app_root"]+"_sim", &simApp)
	return devApp.Status == "up" || simApp.Status == "up"
}
