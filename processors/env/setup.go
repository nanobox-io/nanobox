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

	// if there is no boxfile display a message and end
	if box := boxfile.NewFromPath(config.Boxfile()); !box.Valid {
		display.MissingBoxfile()
		return util.ErrorfQuiet("missing or invalid boxfile")
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
