package component

import (
  "encoding/json"
  "fmt"
  
  "github.com/nanobox-io/nanobox/models"
  "github.com/nanobox-io/nanobox/util/boxfile"
)

type updatePayload struct {
  Config map[string]interface{} `json:"config"`
}

// UpdatePayload returns a string for the update hook payload
func UpdatePayload(c *models.Component) (string, error) {
  config, err := boxfile.ComponentConfig(c)
  if err != nil {
    return "", fmt.Errorf("unable to fetch component config: %s", err.Error())
  }

  payload := updatePayload{
    Config: config,
  }

  switch c.Name {
	case "portal", "logvac", "hoarder", "mist":
		payload.Config["token"] = "123"
	}

  j, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to encode update payload: %s", err.Error())
	}

  return string(j), nil
}
