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
func ConfigurePayload(a *models.App, c *models.Component) (string, error) {
  config, err := boxfile.ComponentConfig(c)
  if err != nil {
    return "", fmt.Errorf("unable to fetch component config: %s", err.Error())
  }
  
  payload := configPayload{
		LogvacHost: a.LocalIPs["logvac"],
		MistHost:   a.LocalIPs["mist"],
		MistToken:  "123",
		Platform:   "local",
		Config:     config,
		Member: member{
			LocalIP: c.InternalIP,
			UID:     1,
			Role:    "primary",
		},
		Component: component{
			Name: c.Name,
			UID:  c.Name,
			ID:   c.ID,
		},
		Users: c.Plan.Users,
	}

	switch c.Name {
	case "portal", "logvac", "hoarder", "mist":
		payload.Config["token"] = "123"
	}

	j, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to encode configure payload: %s", err.Error())
	}

	return string(j), nil
}
