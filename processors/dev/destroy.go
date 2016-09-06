package dev

import (
	"fmt"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors/app"
	"github.com/nanobox-io/nanobox/processors/provider"
)

// Destroy ...
func Destroy(appModel *models.App) error {

	// init docker client
	if err := provider.Init(); err != nil {
		return fmt.Errorf("failed to init docker client: %s", err.Error())
	}

	// destroy the app
	if err := app.Destroy(appModel); err != nil {
		return fmt.Errorf("failed to destroy the app: %s", err.Error())
	}

	return nil
}
