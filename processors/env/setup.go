package env

import (
	"fmt"
	
	"github.com/jcelliott/lumber"
	
	"github.com/nanobox-io/nanobox/models"
)

// Setup sets up the provider and the env mounts
func Setup(env *models.Env) error {

	// ensure the env data has been generated
	if err := env.Generate(); err != nil {
		lumber.Error("env:Setup:models:Env:Generate(): %s", err.Error())
		return fmt.Errorf("failed to initialize the env data: %s", err.Error())
	}

	// setup mounts
	if err := Mount(env); err != nil {
		return fmt.Errorf("failed to setup env mounts: %s", err.Error())
	}

	return nil
}
