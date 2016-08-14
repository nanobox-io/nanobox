package component

import (
	"encoding/json"

	"github.com/nanobox-io/nanobox-boxfile"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util"
)

type (

	// Configure ...
	Configure struct {
		App       models.App
		Component models.Component
		boxfile   boxfile.Boxfile
	}

	// member ...
	member struct {
		LocalIP string `json:"local_ip"`
		UID     int    `json:"uid"`
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
		LogvacHost string                     `json:"logvac_host"`
		MistHost   string                     `json:"mist_host"`
		MistToken  string                     `json:"mist_token"`
		Platform   string                     `json:"platform"`
		Config     map[string]interface{}     `json:"config"`
		Member     member                     `json:"member"`
		Component  component                  `json:"component"`
		Users      []models.ComponentPlanUser `json:"users"`
	}

	// startUpdatePayload ...
	startUpdatePayload struct {
		Config map[string]interface{} `json:"config"`
	}
)

//
func (configure Configure) Run() error {

	// short-circuit if the service has already progressed past this point
	if configure.Component.State != "planned" {
		return nil
	}

	if err := configure.loadBoxfile(); err != nil {
		return err
	}

	if err := configure.runUpdate(); err != nil {
		return err
	}

	if err := configure.runConfigure(); err != nil {
		return err
	}

	if err := configure.runStart(); err != nil {
		return err
	}

	if err := configure.persistComponent(); err != nil {
		return err
	}

	return nil
}

// configurePayload ...
func (configure Configure) configurePayload() string {

	// parse the boxfile to fetch the config node
	config := configure.boxfile.Node(configure.Component.Name).Node("config").Parsed

	payload := configPayload{
		LogvacHost: configure.App.LocalIPs["logvac"],
		MistHost:   configure.App.LocalIPs["mist"],
		MistToken:  "123",
		Platform:   "local",
		Config:     config,
		Member: member{
			LocalIP: configure.Component.InternalIP,
			UID:     1,
			Role:    "primary",
		},
		Component: component{
			Name: configure.Component.Name,
			UID:  configure.Component.Name,
			ID:   configure.Component.ID,
		},
		Users: configure.Component.Plan.Users,
	}

	switch configure.Component.Name {
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
func (configure Configure) startUpdatePayload() string {

	// parse the boxfile to fetch the config node
	boxConfig := configure.boxfile.Node(configure.Component.Name).Node("config")

	payload := startUpdatePayload{boxConfig.Parsed}

	switch configure.Component.Name {
	case PORTAL, LOGVAC, HOARDER, MIST:
		payload.Config["token"] = "123"
	}

	j, err := json.Marshal(payload)
	if err != nil {
		return "{}"
	}

	return string(j)
}

// loadBoxfile loads the new build boxfile from the database
func (configure *Configure) loadBoxfile() error {

	env, err := models.FindEnvByID(configure.App.EnvID)
	configure.boxfile = boxfile.New([]byte(env.BuiltBoxfile))

	return err
}

// runUpdate will run the update hook in the container
func (configure *Configure) runUpdate() error {

	// run update
	_, err := util.Exec(configure.Component.ID, "update", configure.startUpdatePayload(), nil)

	return err
}

// runConfigure will run the configure hook in the container
func (configure *Configure) runConfigure() error {
	// run configure
	_, err := util.Exec(configure.Component.ID, "configure", configure.configurePayload(), nil)

	return err
}

// runStart will run the configure hook in the container
func (configure *Configure) runStart() error {

	// run start
	_, err := util.Exec(configure.Component.ID, "start", configure.startUpdatePayload(), nil)

	return err
}

// persistComponent saves the service entry to the database
func (configure *Configure) persistComponent() error {

	configure.Component.State = ACTIVE
	return configure.Component.Save()
}
