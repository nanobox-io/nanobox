package app

import (
	"fmt"

	"github.com/jcelliott/lumber"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors/component"
	"github.com/nanobox-io/nanobox/util/display"
	"github.com/nanobox-io/nanobox/util/locker"
)

// Start will start all services associated with an app
func Start(appModel *models.App) error {
	locker.LocalLock()
	defer locker.LocalUnlock()

	display.OpenContext("Starting %s", appModel.Name)
	defer display.CloseContext()

	// start all the app components
	if err := component.StartAll(appModel); err != nil {
		return fmt.Errorf("failed to start app components: %s", err.Error())
	}

	// set the status to up
	appModel.Status = "up"
	if err := appModel.Save(); err != nil {
		lumber.Error("app:Start:models.App.Save(): %s", err.Error())
		return fmt.Errorf("failed to persist app status: %s", err.Error())
	}

	return nil
}
