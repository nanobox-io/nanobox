package service

import (
	"encoding/json"
	"errors"

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
	me := models.Service{}
	data.Get(config.AppName(), serviceConfigure.control.Meta["name"], &me)

	logvac := models.Service{}
	data.Get(config.AppName(), LOGVAC, &logvac)

	box := boxfile.New([]byte(serviceConfigure.boxfile.Data))
	boxConfig := box.Node(serviceConfigure.control.Meta["name"]).Node("config")

	pload := configPayload{
		LogvacHost: logvac.InternalIP,
		Platform:   "local",
		Config:     boxConfig.Parsed,
		Member: member{
			LocalIP: me.InternalIP,
			UID:     "1",
			Role:    "primary",
		},
		Component: component{
			Name: "whydoesthismatter",
			UID:  serviceConfigure.control.Meta["name"],
			ID:   me.ID,
		},
		Users: me.Plan.Users,
	}
	if pload.Users == nil {
		pload.Users = []models.User{}
	}
	switch serviceConfigure.control.Meta["name"] {
	case PORTAL, LOGVAC, HOARDER, MIST:
		pload.Config["token"] = "123"
	}
	j, err := json.Marshal(pload)
	if err != nil {
		return "{}"
	}
	return string(j)
}

// startUpdatePayload ...
func (serviceConfigure processServiceConfigure) startUpdatePayload() string {
	boxfile := boxfile.New([]byte(serviceConfigure.control.Meta["boxfile"]))
	boxConfig := boxfile.Node(serviceConfigure.control.Meta["name"]).Node("config")

	pload := startUpdatePayload{boxConfig.Parsed}
	switch serviceConfigure.control.Meta["name"] {
	case PORTAL, LOGVAC, HOARDER, MIST:
		pload.Config["token"] = "123"
	}
	j, err := json.Marshal(pload)
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

// loadService loads the service from the database
func (serviceConfigure *processServiceConfigure) loadService() error {
	// get the service from the database; an error means we could not start a service
	// that wasnt setup (ie saved in the database)
	if err := data.Get(config.AppName(), serviceConfigure.control.Meta["name"], &serviceConfigure.service); err != nil {
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

	if err := data.Put(config.AppName(), serviceConfigure.control.Meta["name"], &serviceConfigure.service); err != nil {
		return err
	}

	return nil
}
