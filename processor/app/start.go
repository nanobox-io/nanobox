package app

import (
	"github.com/jcelliott/lumber"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processor/component"
	"github.com/nanobox-io/nanobox/util/display"
	"github.com/nanobox-io/nanobox/util/locker"
)

// Start start all services associated with an app
// will also destroy any running dev containers
type Start struct {
	App models.App
}

//
func (appStart *Start) Run() error {
	display.StartTask("starting existing components")
	// local lock so no starts or stops can run on this app while I am
	locker.LocalLock()
	defer locker.LocalUnlock()

	// start all the apps services
	componentStartAll := component.StartAll{App: appStart.App}
	if err := componentStartAll.Run(); err != nil {
		display.ErrorTask()
		return err
	}

	display.StopTask()

	// set the app status to up
	err := appStart.upApp()
	return err
}

// upApp sets the app status to up
func (appStart *Start) upApp() error {
	appStart.App.Status = UP
	if err := appStart.App.Save(); err != nil {
		lumber.Error("app:Start:App.Save(): %s", err.Error())
	}

	return nil
}
