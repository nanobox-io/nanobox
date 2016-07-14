package app

import (
	"fmt"

	"github.com/nanobox-io/golang-docker-client"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/data"
	"github.com/nanobox-io/nanobox/util/locker"
)

// processAppStop stop all services associated with an app
// will also destroy any running dev containers
type processAppStop struct {
	control processor.ProcessControl
	app     models.App
}

//
func init() {
	processor.Register("app_stop", appStopFn)
}

//
func appStopFn(control processor.ProcessControl) (processor.Processor, error) {
	appStop := &processAppStop{control: control}
	return appStop, appStop.validateMeta()
}

//
func (appStop *processAppStop) Results() processor.ProcessControl {
	return appStop.control
}

//
func (appStop *processAppStop) Process() error {

	// local lock so no starts or stops can run on this app while I am
	locker.LocalLock()
	defer locker.LocalUnlock()

	// load the application
	if err := appStop.loadApp(); err != nil {
		return err
	}

	// short-circuit if the app isn't up
	if !appStop.isUp() {
		return nil
	}

	// in the service package 'name' represents the name of the service
	// so we need to add a app_name to its control
	appStop.control.Meta["app_name"] = appStop.control.Meta["name"]

	// remove any _dev containers that may be running
	// errors are intentionally not caught because
	// if the container doesnt exist we cant remove it
	appStop.removeDev()

	// stop all services
	if err := processor.Run("service_stop_all", appStop.control); err != nil {
		return err
	}

	// set the app status to down
	if err := appStop.downApp(); err != nil {
		return err
	}

	// unmount the environment
	appStop.control.Meta["directory"] = appStop.app.Directory
	return processor.Run("env_unmount", appStop.control)
}

// loadApp loads the app from the db
func (appStop *processAppStop) loadApp() error {
	return data.Get("apps", appStop.control.Meta["name"], &appStop.app)
}

// remove the development container if one exists
// if not dont complain
func (appStop *processAppStop) removeDev() {
	name := fmt.Sprintf("nanobox_%s", appStop.control.Meta["name"])

	docker.ContainerRemove(name)
}

// downApp sets the app status to down
func (appStop *processAppStop) downApp() error {
	appStop.app.Status = DOWN
	return data.Put("apps", appStop.control.Meta["name"], appStop.app)
}

// validateMeta validates the meta data
// it also sets a default for the name of the app
func (appStop *processAppStop) validateMeta() error {

	// set the name of the app if we are not given one
	if appStop.control.Meta["name"] == "" {
		appStop.control.Meta["name"] = fmt.Sprintf("%s_%s", config.AppName(), appStop.control.Env)
	}

	return nil
}

// the app is concidered up if its status is up
// or if any of its containers are up and running
func (appStop *processAppStop) isUp() bool {
	// if the app says its up.. its up
	if appStop.app.Status == UP {
		return true
	}

	// if any of the apps services are running
	// the app is concidered running
	return appServicesRunning(appStop.control.Meta["name"])
}
