package service

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/nanobox-io/nanobox-boxfile"
	"github.com/nanobox-io/nanobox-golang-stylish"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/data"
)

type (

	// processServiceConfigure ...
	processServiceConfigure struct {
		control processor.ProcessControl
		app     models.App
		service models.Service
		boxfile models.Boxfile
	}

	// member ...
	member struct {
		LocalIP string `json:"local_ip"`
		UID     string `json:"uid"`
		Role    string `json:"role"`
	}

	// component ...
	component struct {
		Name string `json:"name"`
		UID  string `json:"uid"`
		ID   string `json:"id"`
	}

	// configPayload ...
	configPayload struct {
		LogvacHost string                 `json:"logvac_host"`
		MistHost   string                 `json:"mist_host"`
		Platform   string                 `json:"platform"`
		Config     map[string]interface{} `json:"config"`
		Member     member                 `json:"member"`
		Component  component              `json:"component"`
		Users      []models.User          `json:"users"`
	}

	// startUpdatePayload ...
	startUpdatePayload struct {
		Config map[string]interface{} `json:"config"`
	}
)

//
func init() {
	processor.Register("service_configure", serviceConfigureFn)
}

// create a service configure and validate the meta
func serviceConfigureFn(control processor.ProcessControl) (processor.Processor, error) {
	serviceConfigure := processServiceConfigure{control: control}
	return serviceConfigure, serviceConfigure.validateMeta()
}

//
func (serviceConfigure processServiceConfigure) Results() processor.ProcessControl {
	return serviceConfigure.control
}

//
func (serviceConfigure processServiceConfigure) Process() error {

	if err := serviceConfigure.loadApp(); err != nil {
		return err
	}

	if err := serviceConfigure.loadService(); err != nil {
		return err
	}

	// short-circuit if the service has already progressed past this point
	if serviceConfigure.service.State != "planned" {
		return nil
	}

	if err := serviceConfigure.loadBoxfile(); err != nil {
		return err
	}

	if err := serviceConfigure.runUpdate(); err != nil {
		return err
	}

	if err := serviceConfigure.runConfigure(); err != nil {
		return err
	}

	if err := serviceConfigure.runStart(); err != nil {
		return err
	}

	if err := serviceConfigure.persistService(); err != nil {
		return err
	}

	return nil
}

// configurePayload ...
func (serviceConfigure processServiceConfigure) configurePayload() string {

	// create some variables for convenience
	name := serviceConfigure.control.Meta["name"]
	app := serviceConfigure.app
	service := serviceConfigure.service

	// parse the boxfile to fetch the config node
	box := boxfile.New([]byte(serviceConfigure.boxfile.Data))
	config := box.Node(name).Node("config").Parsed

	payload := configPayload{
		LogvacHost: app.LocalIPs["logvac"],
		MistHost:   app.LocalIPs["mist"],
		Platform:   "local",
		Config:     config,
		Member: member{
			LocalIP: service.InternalIP,
			UID:     "1",
			Role:    "primary",
		},
		Component: component{
			Name: name,
			UID:  name,
			ID:   service.ID,
		},
		Users: service.Plan.Users,
	}

	if payload.Users == nil {
		payload.Users = []models.User{}
	}

	switch name {
	case PORTAL, LOGVAC, HOARDER, MIST:
		payload.Config["token"] = "123"
	}

	j, err := json.Marshal(payload)
	if err != nil {
		return "{}"
	}

	return string(j)
}

// startUpdatePayload ...
func (serviceConfigure processServiceConfigure) startUpdatePayload() string {

	// parse the boxfile to fetch the config node
	boxfile := boxfile.New([]byte(serviceConfigure.control.Meta["boxfile"]))
	boxConfig := boxfile.Node(serviceConfigure.control.Meta["name"]).Node("config")

	payload := startUpdatePayload{boxConfig.Parsed}

	switch serviceConfigure.control.Meta["name"] {
	case PORTAL, LOGVAC, HOARDER, MIST:
		payload.Config["token"] = "123"
	}

	j, err := json.Marshal(payload)
	if err != nil {
		return "{}"
	}

	return string(j)
}

// validateMeta validates that the image is provided
func (serviceConfigure *processServiceConfigure) validateMeta() error {
	// make sure i was given a name and image
	if serviceConfigure.control.Meta["name"] == "" {
		return errors.New("missing name")
	}

	return nil
}

// loadApp loads the app from the database
func (serviceConfigure *processServiceConfigure) loadApp() error {

	// load the app from the database
	key := fmt.Sprintf("%s_%s", config.AppName(), serviceConfigure.control.Env)
	if err := data.Get("apps", key, &serviceConfigure.app); err != nil {
		return err
	}

	return nil
}

// loadService loads the service from the database
func (serviceConfigure *processServiceConfigure) loadService() error {

	name := serviceConfigure.control.Meta["name"]

	// get the service from the database; an error means we could not start a service
	// that wasnt setup (ie saved in the database)
	bucket := fmt.Sprintf("%s_%s", config.AppName(), serviceConfigure.control.Env)
	if err := data.Get(bucket, name, &serviceConfigure.service); err != nil {
		return err
	}

	return nil
}

// loadBoxfile loads the new build boxfile from the database
func (serviceConfigure *processServiceConfigure) loadBoxfile() error {

	// we won't worry about erroring here, because there may not be
	// a build_boxfile at this point
	data.Get(config.AppName()+"_meta", "build_boxfile", &serviceConfigure.boxfile)

	return nil
}

// runUpdate will run the update hook in the container
func (serviceConfigure *processServiceConfigure) runUpdate() error {
	// run update
	serviceConfigure.control.Info(stylish.SubBullet("Updating services..."))
	_, err := util.Exec(serviceConfigure.service.ID, "update", serviceConfigure.startUpdatePayload(), nil)

	return err
}

// runConfigure will run the configure hook in the container
func (serviceConfigure *processServiceConfigure) runConfigure() error {
	// run configure
	serviceConfigure.control.Info(stylish.SubBullet("Configuring services..."))
	_, err := util.Exec(serviceConfigure.service.ID, "configure", serviceConfigure.configurePayload(), nil)

	return err
}

// runStart will run the configure hook in the container
func (serviceConfigure *processServiceConfigure) runStart() error {
	// run start
	serviceConfigure.control.Info(stylish.SubBullet("Starting services..."))
	_, err := util.Exec(serviceConfigure.service.ID, "start", serviceConfigure.startUpdatePayload(), nil)

	return err
}

// persistService saves the service entry to the database
func (serviceConfigure *processServiceConfigure) persistService() error {
	serviceConfigure.service.State = ACTIVE

	name := serviceConfigure.control.Meta["name"]

	bucket := fmt.Sprintf("%s_%s", config.AppName(), serviceConfigure.control.Env)
	if err := data.Put(bucket, name, &serviceConfigure.service); err != nil {
		return err
	}

	return nil
}
