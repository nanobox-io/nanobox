package component

import (
	"encoding/json"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util/boxfile"
)

// PlanPayload returns a string for the user hook payload
func PlanPayload(component *models.Component) string {
	config, err := boxfile.ComponentConfig(component)
	if err != nil {
		return "{}"
	}

	payload := map[string]interface{}{
		"config": config,
	}

	// marshal the payload into json
	b, err := json.Marshal(payload)
	if err != nil {
		return "{}"
	}

	return string(b)
}
