package dev

import (
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processor/app"
)

// Destroy ...
type Destroy struct {
	App models.App
}

//
func (destroy Destroy) Run() error {
	appDestroy := app.Destroy{App: destroy.App}

	return appDestroy.Run()
}
