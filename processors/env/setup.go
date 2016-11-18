package env

import (
	"fmt"

	"github.com/jcelliott/lumber"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors/provider"
	"github.com/nanobox-io/nanobox/util/display"
	util_provider "github.com/nanobox-io/nanobox/util/provider"
)

// Setup sets up the provider and the env mounts
func Setup(envModel *models.Env) error {

	// init docker client
	if err := provider.Init(); err != nil {
		return fmt.Errorf("failed to init docker client: %s", err.Error())
	}

	// ensure the envModel data has been generated
	if err := envModel.Generate(); err != nil {
		lumber.Error("env:Setup:models:Env:Generate(): %s", err.Error())
		return fmt.Errorf("failed to initialize the env data: %s", err.Error())
	}

	if util_provider.HasMount(fmt.Sprintf("%s%s/code", util_provider.HostShareDir(), envModel.ID)) {
		return nil
	}

	display.OpenContext("Preparing environment")
	defer display.CloseContext()

	// setup mounts
	if err := Mount(envModel); err != nil {
		display.ErrorTask()
		return fmt.Errorf("failed to setup env mounts: %s", err.Error())
	}

	return nil
}
