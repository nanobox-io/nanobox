package dev

import (
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors/app"
)

// Destroy ...
func Destroy(appModel *models.App) error {

	return app.Destroy(appModel)
}
