package dev

import (
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors/component"
)

// Deploy ...
func Deploy(envModel *models.Env, appModel *models.App) error {

	// syncronize the services as per the new boxfile
	if err := component.Sync(envModel, appModel); err != nil {
		return err
	}

	return nil
}
