package env

import (
	"fmt"

	"github.com/jcelliott/lumber"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util/display"
)

// Setup sets up the provider and the env mounts
func Setup(envModel *models.Env) error {

	display.OpenContext("setting up environment")
	defer display.CloseContext()

	// ensure the envModel data has been generated
	if err := envModel.Generate(); err != nil {
		lumber.Error("env:Setup:models:Env:Generate(): %s", err.Error())
		return fmt.Errorf("failed to initialize the env data: %s", err.Error())
	}

	// setup mounts
	if err := Mount(envModel); err != nil {
		display.ErrorTask()
		return fmt.Errorf("failed to setup env mounts: %s", err.Error())
	}

	return nil
}
