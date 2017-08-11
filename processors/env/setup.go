package env

import (
	"fmt"
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

	// init docker client
	if err := provider.Init(); err != nil {
		return util.ErrorAppend(err, "failed to init docker client")
	}

	// ensure the envModel data has been generated
	if err := envModel.Generate(); err != nil {
		lumber.Error("env:Setup:models:Env:Generate()")
		return util.ErrorAppend(err, "failed to initialize the env data")
	}

	if util_provider.HasMount(fmt.Sprintf("%s%s/code", util_provider.HostShareDir(), envModel.ID)) {
		return nil
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
