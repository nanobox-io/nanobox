package app

import (
	"github.com/jcelliott/lumber"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors/component"
	"github.com/nanobox-io/nanobox/util"
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
			return util.ErrorAppend(err, "failed to setup the app")
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
		return util.ErrorAppend(err, "failed to clean crufty components")
	}

	// start all the app components
	if err := component.StartAll(appModel); err != nil {
		return util.ErrorAppend(err, "failed to start app components")
	}

	// set the status to up
	appModel.Status = "up"
	if err := appModel.Save(); err != nil {
		lumber.Error("app:Start:models.App.Save()")
		return util.ErrorAppend(err, "failed to persist app status")
	}

	return nil
}
