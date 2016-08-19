package sim

import (

	"github.com/jcelliott/lumber"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors/app"
)

// Destroy ...
type Destroy struct {
	App models.App
}

//
func (destroy *Destroy) Run() error {
	lumber.Debug("simDestroy:App: %+v", destroy.App)
	appDestroy := app.Destroy{App: destroy.App}

	return appDestroy.Run()
}
