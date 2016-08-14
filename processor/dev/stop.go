package dev

import (
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processor/app"
)

// Stop ...
type Stop struct {
	App models.App
}

//
func (stop Stop) Run() error {
	appStop := app.Stop{
		App: stop.App,
	}
	return appStop.Run()
}
