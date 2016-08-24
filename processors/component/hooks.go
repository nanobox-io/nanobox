package component

import (
  "fmt"
  
  "github.com/jcelliott/lumber"
  
  "github.com/nanobox-io/nanobox/generators/hooks/component"
  "github.com/nanobox-io/nanobox/models"
  "github.com/nanobox-io/nanobox/util/hookit"
)

// RunPlanHook runs the plan hook inside of the specified container
func RunPlanHook(c *models.Component) (string, error) {
  // generate the plan payload
  planPayload, err := component.PlanPayload(c)
  if err != nil {
    lumber.Error("component:RunPlanHook:component.PlanPayload(): %s", err.Error())
    return "", fmt.Errorf("failed to generate payload for plan hook: %s", err.Error())
  }
  
  // run the plan hook
  res, err := hookit.Exec(c.ID, "plan", planPayload, "debug")
  if err != nil {
    lumber.Error("component:RunPlanHook:hookit.Exec(%s, %s, %s, %s): %s", c.ID, "plan", planPayload, "debug", err.Error())
    return "", fmt.Errorf("failed to execute plan hook: %s", err.Error())
  }
  
  return res, nil
}

// RunConfigureHook runs the configure hook inside of the specified container
func RunConfigureHook(a *models.App, c *models.Component) (string, error) {
  // generate the configure payload
  configurePayload, err := component.ConfigurePayload(a, c)
  if err != nil {
    lumber.Error("component:RunConfigureHook:component.ConfigurePayload(): %s", err.Error())
    return "", fmt.Errorf("failed to generate payload for configure hook: %s", err.Error())
  }
  
  // run the configure hook
  res, err := hookit.Exec(c.ID, "configure", configurePayload, "debug")
  if err != nil {
    lumber.Error("component:RunConfigureHook:hookit.Exec(%s, %s, %s, %s): %s", c.ID, "configure", configurePayload, "debug", err.Error())
    return "", fmt.Errorf("failed to execute configure hook: %s", err.Error())
  }
  
  return res, nil
}

// RunStartHook runs the start hook inside of the specified container
func RunStartHook(c *models.Component) (string, error) {
  // generate the start payload
  startPayload, err := component.StartPayload(c)
  if err != nil {
    lumber.Error("component:RunStartHook:component.StartPayload(): %s", err.Error())
    return "", fmt.Errorf("failed to generate payload for start hook: %s", err.Error())
  }
  
  // run the start hook
  res, err := hookit.Exec(c.ID, "start", startPayload, "debug")
  if err != nil {
    lumber.Error("component:RunStartHook:hookit.Exec(%s, %s, %s, %s): %s", c.ID, "start", startPayload, "debug", err.Error())
    return "", fmt.Errorf("failed to execute start hook: %s", err.Error())
  }
  
  return res, nil
}

// RunUpdateHook runs the update hook inside of the specified container
func RunUpdateHook(c *models.Component) (string, error) {
  // generate the update payload
  updatePayload, err := component.UpdatePayload(c)
  if err != nil {
    lumber.Error("component:RunUpdateHook:component.UpdatePayload(): %s", err.Error())
    return "", fmt.Errorf("failed to generate payload for update hook: %s", err.Error())
  }
  
  // run the update hook
  res, err := hookit.Exec(c.ID, "update", updatePayload, "debug")
  if err != nil {
    lumber.Error("component:RunUpdateHook:hookit.Exec(%s, %s, %s, %s): %s", c.ID, "update", updatePayload, "debug", err.Error())
    return "", fmt.Errorf("failed to execute update hook: %s", err.Error())
  }
  
  return res, nil
}
