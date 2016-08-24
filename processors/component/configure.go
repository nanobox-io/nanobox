package component

import (
  "fmt"
  
  "github.com/jcelliott/lumber"
  
  "github.com/nanobox-io/nanobox/models"
)

// Configure configures a component
func Configure(a *models.App, c *models.Component) error {
  
  // short-circuit if the component state is not planned
  if c.State != "planned" {
    return nil
  }
  
  // run the update hook
  if _, err := RunUpdateHook(c); err != nil {
    return fmt.Errorf("failed to run update hook: %s", err.Error())
  }
  
  // run the configure hook
  if _, err := RunConfigureHook(a, c); err != nil {
    return fmt.Errorf("failed to run configure hook: %s", err.Error())
  }
  
  // run the start hook
  if _, err := RunStartHook(c); err != nil {
    return fmt.Errorf("failed to run start hook: %s", err.Error())
  }
  
  // set state as active
  c.State = "active"
  if err := c.Save(); err != nil {
    lumber.Error("component:Configure:models.Component.Save()", err.Error())
    return fmt.Errorf("failed to set component state: %s", err.Error())
  }
  
  return nil
}
