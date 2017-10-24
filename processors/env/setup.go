package env

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/jcelliott/lumber"
	"github.com/nanobox-io/nanobox-boxfile"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors/provider"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/display"
	util_provider "github.com/nanobox-io/nanobox/util/provider"
)

// Setup sets up the provider and the env mounts
func Setup(envModel *models.Env) error {
	lumber.Trace("env:Setup:config:Boxfile() - %s", config.Boxfile())

	// check for boxfile in separate step from valid yaml syntax
	// if there is no boxfile, display a message and end
	box, err := boxfile.NewFromFile(config.Boxfile())
	if err != nil || box == nil {
		// todo: recursively check for boxfile
		display.MissingBoxfile()
		err = util.ErrorfQuiet("[USER] missing boxfile - %s", err.Error())
		if err2, ok := err.(util.Err); ok {
			err2.Suggest = "Ensure you have a boxfile.yml file in your current app directory"
			return err2
		}
		return err
	}

	// if boxfile is invalid, display a message and end
	if !box.Valid {
		display.InvalidBoxfile()
		err = util.ErrorfQuiet("[USER] invalid yaml found in boxfile")
		if err2, ok := err.(util.Err); ok {
			err2.Suggest = "It appears you have an invalid boxfile. Validate it at `yamllint.com`"
			return err2
		}
		return err
	}

	// todo: validate boxfile nodes

	// ensure local engine exists
	engineDir, err := config.EngineDir()
	if err != nil {
		display.LocalEngineNotFound()
		err = util.ErrorfQuiet("[USER] custom local engine not found - %s", err.Error())
		if err2, ok := err.(util.Err); ok {
			err2.Suggest = "Ensure the engine defined in your boxfile.yml exists at the location specified"
			return err2
		}
		return err
	}

	// init docker client
	if err := provider.Init(); err != nil {
		return util.ErrorAppend(err, "failed to init docker client")
	}

	// ensure the envModel data has been generated
	if err := envModel.Generate(); err != nil {
		lumber.Error("env:Setup:models:Env:Generate()")
		return util.ErrorAppend(err, "failed to initialize the env data")
	}

	// if switch from local engine, ensure old local engine gets unmounted
	oldBox := boxfile.New([]byte(envModel.UserBoxfile))
	newBox, _ := boxfile.NewFromFile(config.Boxfile()) // can ignore error, we made sure it exists before this point
	oldEngineName := oldBox.Node("run.config").StringValue("engine")
	newEngineName := newBox.Node("run.config").StringValue("engine")
	if (oldEngineName != newEngineName) && newEngineName != "" {
		oldEnginePath, err := filepath.Abs(oldEngineName)
		if err != nil {
			// todo: ignore here so if they delete their engine it won't break until they restore it
			return fmt.Errorf("Failed to resolve old engine location - %s", err.Error())
		}

		err = UnmountEngine(envModel, oldEnginePath)
		if err != nil {
			return fmt.Errorf("Failed to unmount engine - %s", err.Error())
		}
	}

	if util_provider.HasMount(fmt.Sprintf("%s%s/code", util_provider.HostShareDir(), envModel.ID)) {
		if engineDir == "" {
			return nil
		}
		if util_provider.HasMount(fmt.Sprintf("%s%s/engine", util_provider.HostShareDir(), envModel.ID)) {
			return nil
		}
	}

	display.OpenContext("Preparing environment")
	defer display.CloseContext()

	// setup mounts
	mountFunc := func() error {
		return Mount(envModel)
	}

	if err := util.Retry(mountFunc, 5, (time.Second * 10)); err != nil {
		display.ErrorTask()
		return util.ErrorAppend(err, "failed to setup env mounts")
	}

	return nil
}
