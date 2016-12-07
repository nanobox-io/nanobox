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
func Start(envModel *models.Env, appModel *models.App, name string) error {

	display.OpenContext("%s (%s)", envModel.Name, appModel.DisplayName())
	defer display.CloseContext()

	// if the app been initialized run the setup
	if appModel.State != "active" {
		if err := Setup(envModel, appModel, name); err != nil {
			return fmt.Errorf("failed to setup the app: %s", err)
		}
	} else {
		// restoring app
		display.StartTask("Restoring App")
		display.StopTask()
	}

	// we reserver here only while people are transitioning
	// this can go away once everyone is on the new natless method
	reserveIPs(appModel)

	locker.LocalLock()
	defer locker.LocalUnlock()

	// clean crufty components
	if err := component.Clean(appModel); err != nil {
		return fmt.Errorf("failed to clean crufty components: %s", err.Error())
	}

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
