package platform

import (
  "fmt"
  
  "github.com/jcelliott/lumber"
  
  "github.com/nanobox-io/nanobox/models"
  "github.com/nanobox-io/nanobox/processors/component"
)

// provisionComponent will provision an individual component
func provisionComponent(a *models.App, pc PlatformComponent) error {
  
  // if the component exists and is active just start it and stop here
  if isComponentActive(a, pc.name) {
    c, _ := models.FindComponentBySlug(a.ID, pc.name)
    
    // start the component
    if err := component.Start(c); err != nil {
      return fmt.Errorf("failed to start component: %s", err.Error())
    }
    
    return nil
  }
  
  // setup
  if err := component.Setup(a, pc.name, pc.label, pc.image); err != nil {
    return fmt.Errorf("failed to setup platform component (%s): %s", pc.label, err.Error())
  }
  
  // load the component
  c, err := models.FindComponentBySlug(a.ID, pc.name)
  if err != nil {
    lumber.Error("platform:provisionComponent:models.FindComponentBySlug(%s, %s): %s", a.ID, pc.name, err.Error())
    return fmt.Errorf("failed to load component from db: %s", err.Error())
  }
  
  // configure
  if err := component.Configure(a, c); err != nil {
    return fmt.Errorf("failed to configure platform component (%s): %s", pc.label, err.Error())
  }
  
  return nil
}

// isComponentActive returns true if a component is already active
func isComponentActive(a *models.App, name string) bool {
	// component db entry
	component, _ := models.FindComponentBySlug(a.ID, name)

	return component.State == "active"
}
