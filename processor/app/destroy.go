package app

import (
	"fmt"
	"net"
	"strings"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/data"
	"github.com/nanobox-io/nanobox/util/dhcp"
	"github.com/nanobox-io/nanobox/util/locker"
)

// processAppDestroy ...
type processAppDestroy struct {
	control processor.ProcessControl
	app     models.App
}

//
func init() {
	processor.Register("app_destroy", appDestroyFn)
}

//
func appDestroyFn(control processor.ProcessControl) (processor.Processor, error) {
	appDestroy := &processAppDestroy{control: control}
	return appDestroy, appDestroy.validateMeta()
}

func (appDestroy *processAppDestroy) validateMeta() error {
	if appDestroy.control.Env == "" {
		return fmt.Errorf("Env not set")
	}

	if appDestroy.control.Meta["name"] == "" {
		appDestroy.control.Meta["name"] = fmt.Sprintf("%s_%s", config.AppID(), appDestroy.control.Env)
	}

	return nil
}

//
func (appDestroy *processAppDestroy) Results() processor.ProcessControl {
	return appDestroy.control
}

//
func (appDestroy *processAppDestroy) Process() error {

	// establish an app-level lock to ensure we're the only ones setting up an app
	// also, we need to ensure that the lock is released even if we error out.
	locker.LocalLock()
	defer locker.LocalUnlock()

	if err := appDestroy.loadApp(); err != nil {
		return err
	}

	if err := appDestroy.releaseIPs(); err != nil {
		return err
	}

	if err := appDestroy.deleteMeta(); err != nil {
		return err
	}

	if err := appDestroy.deleteApp(); err != nil {
		return err
	}

	return nil
}

// loadApp loads the app from the db
func (appDestroy *processAppDestroy) loadApp() error {

	// durring the app destroy it should absolutely exist
	// so we will be returning an error if it fails
	return data.Get("apps", appDestroy.control.Meta["name"], &appDestroy.app)
}

// releaseIPs releases necessary app-global ip addresses
func (appDestroy *processAppDestroy) releaseIPs() error {

	// release all of the external IPs
	for _, ip := range appDestroy.app.GlobalIPs {
		// release the IP
		if err := dhcp.ReturnIP(net.ParseIP(ip)); err != nil {
			return err
		}
	}

	// release all of the local IPs
	for _, ip := range appDestroy.app.LocalIPs {
		// release the IP
		if err := dhcp.ReturnIP(net.ParseIP(ip)); err != nil {
			return err
		}
	}

	return nil
}

// deleteMeta deletes metadata about this app from the database
func (appDestroy *processAppDestroy) deleteMeta() error {
	// just get the raw appid without the dev or sim
	stripped := strings.Replace(appDestroy.control.Meta["name"], "_dev", "", -1)
	stripped = strings.Replace(stripped, "_sim", "", -1)

	// remove the boxfile information
	data.Delete(stripped+"_meta", "build_boxfile")
	data.Delete(stripped+"_meta", "dev_build_boxfile")

	// delete the evars model
	return data.Delete(stripped+"_meta", appDestroy.control.Env+"_env")
}

// deleteApp deletes the app to the db
func (appDestroy *processAppDestroy) deleteApp() error {

	// delete the app model
	if err := data.Delete("apps", appDestroy.control.Meta["name"]); err != nil {
		return err
	}

	return nil
}
