package platform

import (
  "fmt"
  
  "github.com/jcelliott/lumber"
  
  "github.com/nanobox-io/nanobox/models"
  "github.com/nanobox-io/nanobox/processors/component"
)

// Stop stops all platform components
func Stop(a *models.App) error {
  for _, pc := range append(setupComponents, deployComponents...) {
    if err := stopComponent(a, pc); err != nil {
      return fmt.Errorf("failed to stop platform component: %s", err.Error())
    }
  }
  
  return nil
}

// stopComponent stops a platform component
func stopComponent(a *models.App, pc PlatformComponent) error {
  // load the component
  c, err := models.FindComponentBySlug(a.ID, pc.name)
  if err != nil {
    lumber.Error("platform:stopComponent:models.FindComponentBySlug(%s, %s): %s", a.ID, pc.name, err.Error())
    return fmt.Errorf("failed to load component: %s", err.Error())
  }
  
  // stop the component
  if err := component.Stop(c); err != nil {
    return fmt.Errorf("failed to stop component: %s", err.Error())
  }
  
  return nil
}
