package sim

import (
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors/app"
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
