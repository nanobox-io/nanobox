package dev

import (
	"fmt"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors/component"
	"github.com/nanobox-io/nanobox/processors/env"
)

// Deploy ...
func Deploy(envModel *models.Env, appModel *models.App) error {

	// init docker client
	if err := env.Setup(envModel); err != nil {
		return fmt.Errorf("failed to init docker client: %s", err.Error())
	}

	// syncronize the services as per the new boxfile
	if err := component.Sync(envModel, appModel); err != nil {
		return fmt.Errorf("failed to sync components: %s", err.Error())
	}

	return nil
}
