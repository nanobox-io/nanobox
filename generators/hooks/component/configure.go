package component

import (
	"encoding/json"

	"github.com/nanobox-io/nanobox/models"
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
func ConfigurePayload(appModel *models.App, componentModel *models.Component) string {
	config, err := componentConfig(componentModel)
	if err != nil {
		// lumber.Error("unable to fetch component config: %s", err.Error())
		return "{}"
	}

	payload := configPayload{
		LogvacHost: appModel.LocalIPs["logvac"],
		MistHost:   appModel.LocalIPs["mist"],
		MistToken:  "123",
		Platform:   "local",
		Config:     config,
		Member: member{
			LocalIP: componentModel.IPAddr(),
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
		return "{}"
	}

	return string(j)
}
