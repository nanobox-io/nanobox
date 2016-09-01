package dev

import (
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors/component"
	"github.com/nanobox-io/nanobox/processors/provider"
)

// Deploy ...
func Deploy(envModel *models.Env, appModel *models.App) error {

	// init docker client
	if err := provider.Init(); err != nil {
		return fmt.Errorf("failed to init docker client: %s", err.Error())
	}

	// syncronize the services as per the new boxfile
	if err := component.Sync(envModel, appModel); err != nil {
		return err
	}

	return nil
}
