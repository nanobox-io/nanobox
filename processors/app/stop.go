package app

import (
  "fmt"
  
  "github.com/jcelliott/lumber"
  
  "github.com/nanobox-io/nanobox/models"
  "github.com/nanobox-io/nanobox/processors/component"
  "github.com/nanobox-io/nanobox/processors/provider"
)

// Stop will stop all services associated with an app
func Stop(a *models.App) error {
  locker.LocalLock()
  defer locker.LocalUnlock()
  
  // short-circuit if the app is already down
  // TODO: also check if any containers are running
  if a.status != "up" {
    return nil
  }
  
  // initialize docker for the provider
  if err := provider.Init(); err != nil {
    return fmt.Errorf("failed to initialize docker environment: %s", err.Error())
  }
  
  // stop all app components
  if err := component.StopAll(a); err != nil {
    return fmt.Errorf("failed to stop all app components: %s", err.Error())
  }
  
  // stop any dev containers
  docker.ContainerRemove(fmt.Sprintf("nanobox_%s", a.ID))
  
  // set the status to down
  a.status = "down"
  if err := a.Save(); err != nil {
    lumber.Error("app:Stop:models.App.Save(): %s", err.Error())
    return fmt.Errorf("failed to persist app status: %s", err.Error())
  }
  
  return nil
}
