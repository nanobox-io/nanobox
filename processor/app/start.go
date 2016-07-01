package app

import (
	"fmt"
	
	"github.com/nanobox-io/nanobox/util/locker"
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/data"
)

// processAppStart start all services associated with an app
// will also destroy any running dev containers
type processAppStart struct {
	control processor.ProcessControl
	app     models.App
}

//
func init() {
	processor.Register("app_start", appStartFn)
}

//
func appStartFn(control processor.ProcessControl) (processor.Processor, error) {
	appStart := &processAppStart{control: control}
	return appStart, appStart.validateMeta()
}

//
func (appStart *processAppStart) Results() processor.ProcessControl {
	return appStart.control
}

//
func (appStart *processAppStart) Process() error {

	// local lock so no starts or stops can run on this app while I am
	locker.LocalLock()
	defer locker.LocalUnlock()

	// the env setup is idempotent and will not output anything 
	// unless it actually does something
	if err := processor.Run("env_setup", appStart.control); err != nil {
		return err
	}

	if err := appStart.loadApp(); err != nil {
		return err
	}

	// // if the app has not been setup yet we probably need to run 
	// // the env setup
	// if appStart.app.State != ACTIVE || !provider.IsReady() {

	// 	// output message
	// 	appStart.control.Display("it appears the environment wasnt up to run this command")
	// 	appStart.control.Display("We will get that for you.")

	// 	// run env_setup
	// 	if err := processor.Run("env_setup", appStart.control); err != nil {
	// 		return err
	// 	}
	// }

	// in the service package 'name' represents the name of the service
	// so we need to add a app_name to its control
	appStart.control.Meta["app_name"] = appStart.control.Meta["name"]

	// start all the apps services
	if err := processor.Run("service_start_all", appStart.control); err != nil {
		return err
	}

	// set the app status to up
	return appStart.upApp()
}

// appExists checks to see if the app is in the database
func (appStart *processAppStart) appExists() bool {
	app := models.App{}
	err := data.Get("apps", appStart.control.Meta["name"], &app)
	return err == nil
}

// loadApp loads the app from the db
func (appStart *processAppStart) loadApp() error {
	return data.Get("apps", appStart.control.Meta["name"], &appStart.app)
}

// upApp sets the app status to up
func (appStart *processAppStart) upApp() error {
	appStart.app.Status = UP
	return data.Put("apps", appStart.control.Meta["name"], appStart.app)
}

// validateMeta validates the meta data
// it also sets a default for the name of the app
func (appStart *processAppStart) validateMeta() error {

	// set the name of the app if we are not given one
	if appStart.control.Meta["name"] == "" {
		appStart.control.Meta["name"] = fmt.Sprintf("%s_%s", config.AppName(), appStart.control.Env)
	}

	return nil
}
