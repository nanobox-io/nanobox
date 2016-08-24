package component

import (
  "encoding/json"
  "fmt"
  
  "github.com/nanobox-io/nanobox/models"
  "github.com/nanobox-io/nanobox/util/boxfile"
)

// PlanPayload returns a string for the user hook payload
func PlanPayload(component *models.Component) (string, error) {
  config, err := boxfile.ComponentConfig(component)
  if err != nil {
    return "", fmt.Errorf("failed to fetch component config: %s", err.Error())
  }
  
  payload := map[string]interface{}{
    "config": config,
  }
  
  // marshal the payload into json
  b, err := json.Marshal(payload)
  if err != nil {
    return "", fmt.Errorf("failed to encode hook payload into json: %s", err.Error())
  }
  
  return string(b), nil
}
