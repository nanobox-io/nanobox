package component

import (
	"encoding/json"

	"github.com/nanobox-io/nanobox/models"
)

type startPayload struct {
	Config map[string]interface{} `json:"config"`
}

// StartPayload returns a string for the start hook payload
func StartPayload(c *models.Component) string {
	config, err := componentConfig(c)
	if err != nil {
		return "{}"
	}

	payload := startPayload{
		Config: config,
	}

	switch c.Name {
	case "portal", "logvac", "hoarder", "mist":
		payload.Config["token"] = "123"
	}

	j, err := json.Marshal(payload)
	if err != nil {
		return "{}"
	}

	return string(j)
}
