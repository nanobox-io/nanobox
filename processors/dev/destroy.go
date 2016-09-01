package dev

import (
	"fmt"
	
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors/app"
)

// Destroy ...
func Destroy(appModel *models.App) error {

	// destroy the app
	if err := app.Destroy(appModel); err != nil {
		return fmt.Errorf("failed to destroy the app: %s", err.Error())
	}

	return nil
}
