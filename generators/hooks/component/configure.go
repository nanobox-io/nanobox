package component

import (
	"encoding/json"
	"fmt"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util/boxfile"
)

// member ...
type member struct {
	LocalIP string `json:"local_ip"`
	UID     int    `json:"uid"`
	Role    string `json:"role"`
}

// component ...
type component struct {
	Name string `json:"name"`
	UID  string `json:"uid"`
	ID   string `json:"id"`
}

// configPayload ...
type configPayload struct {
	LogvacHost string                     `json:"logvac_host"`
	MistHost   string                     `json:"mist_host"`
	MistToken  string                     `json:"mist_token"`
	Platform   string                     `json:"platform"`
	Config     map[string]interface{}     `json:"config"`
	Member     member                     `json:"member"`
	Component  component                  `json:"component"`
	Users      []models.ComponentPlanUser `json:"users"`
}

// ConfigurePayload returns a string for the configure hook payload
func ConfigurePayload(appModel *models.App, componentModel *models.Component) (string, error) {
	config, err := boxfile.ComponentConfig(componentModel)
	if err != nil {
		return "", fmt.Errorf("unable to fetch component config: %s", err.Error())
	}

	payload := configPayload{
		LogvacHost: appModel.LocalIPs["logvac"],
		MistHost:   appModel.LocalIPs["mist"],
		MistToken:  "123",
		Platform:   "local",
		Config:     config,
		Member: member{
			LocalIP: componentModel.InternalIP,
			UID:     1,
			Role:    "primary",
		},
		Component: component{
			Name: componentModel.Name,
			UID:  componentModel.Name,
			ID:   componentModel.ID,
		},
		Users: componentModel.Plan.Users,
	}

	j, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to encode configure payload: %s", err.Error())
	}

	return string(j), nil
}
