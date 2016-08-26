package dev

import (
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors/app"
)

//
func Stop(appModel *models.App) error {
	// do something dev specific

	return app.Stop(appModel)
}
