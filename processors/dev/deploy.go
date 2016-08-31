package dev

import (
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors/component"
	"github.com/nanobox-io/nanobox/processors/provider"
	"github.com/nanobox-io/nanobox/util/display"
)

// Deploy ...
func Deploy(envModel *models.Env, appModel *models.App) error {
	display.OpenContext("deploying dev")
	defer display.CloseContext()

	// run the share init which gives access to docker
	if err := provider.Init(); err != nil {
		return err
	}

	// syncronize the services as per the new boxfile
	if err := component.Sync(envModel, appModel); err != nil {
		return err
	}

	return nil
}
