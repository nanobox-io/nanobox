package dev

import (
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors/platform"
)

//
func Log(appModel *models.App) error {

	return platform.MistListen(appModel)
}
