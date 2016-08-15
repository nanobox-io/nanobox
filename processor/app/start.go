package app

import (
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processor/component"
	"github.com/nanobox-io/nanobox/util/locker"
)

// Start start all services associated with an app
// will also destroy any running dev containers
type Start struct {
	App models.App
}

//
func (appStart *Start) Run() error {

	// local lock so no starts or stops can run on this app while I am
	locker.LocalLock()
	defer locker.LocalUnlock()

	// start all the apps services
	componentStartAll := component.StartAll{App: appStart.App}
	if err := componentStartAll.Run(); err != nil {
		return err
	}

	// set the app status to up
	return appStart.upApp()
}

// upApp sets the app status to up
func (appStart *Start) upApp() error {
	appStart.App.Status = UP
	return appStart.App.Save()
}
