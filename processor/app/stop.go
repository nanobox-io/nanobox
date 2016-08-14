package app

import (
	"fmt"

	"github.com/nanobox-io/golang-docker-client"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processor/component"
	"github.com/nanobox-io/nanobox/util/locker"
)

// Stop stop all services associated with an app
// will also destroy any running dev containers
type Stop struct {
	App models.App
}

//
func (stop *Stop) Run() error {

	// local lock so no starts or stops can run on this app while I am
	locker.LocalLock()
	defer locker.LocalUnlock()

	// short-circuit if the app isn't up
	if !stop.isUp() {
		return nil
	}

	// intitialize the environment
	dockerInit()

	// remove any _dev containers that may be running
	// errors are intentionally not caught because
	// if the container doesnt exist we cant remove it
	stop.removeDev()

	// stop all services
	componentStopAll := component.StopAll{App: stop.App}
	if err := componentStopAll.Run(); err != nil {
		return err
	}

	// set the app status to down
	if err := stop.downApp(); err != nil {
		return err
	}

	return nil
}

// remove the development container if one exists
// if not dont complain
func (stop *Stop) removeDev() {
	name := fmt.Sprintf("nanobox_%s", stop.App.ID)

	docker.ContainerRemove(name)
}

// downApp sets the app status to down
func (stop *Stop) downApp() error {
	stop.App.Status = DOWN
	return stop.App.Save()
}

// the app is concidered up if its status is up
// or if any of its containers are up and running
func (stop *Stop) isUp() bool {
	// if the app says its up.. its up
	if stop.App.Status == UP {
		return true
	}

	// if any of the apps services are running
	// the app is concidered running
	return appServicesRunning(stop.App)
}
