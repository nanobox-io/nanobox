package dev

import (
	"fmt"
	
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors/component"
	"github.com/nanobox-io/nanobox/processors/env"
	"github.com/nanobox-io/nanobox/processors/provider"
)

// Deploy ...
func Deploy(envModel *models.Env, appModel *models.App) error {

	// ensure the provider is setup
	if err := provider.Setup(); err != nil {
		return fmt.Errorf("failed to setup the provider: %s", err.Error())
	}

	// ensure the environment is setup
	if err := env.Setup(envModel); err != nil {
		return fmt.Errorf("failed to setup the env: %s", err.Error())
	}

	// syncronize the services as per the new boxfile
	if err := component.Sync(envModel, appModel); err != nil {
		return err
	}

	return nil
}
