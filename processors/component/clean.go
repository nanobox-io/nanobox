package component

import (
  "fmt"
  
  "github.com/jcelliott/lumber"
  "github.com/nanobox-io/golang-docker-client"
  
  "github.com/nanobox-io/nanobox/models"
)

// Clean purges any components in a dirty or incomplete state
func Clean(a *models.App) error {
  // fetch all of the app components
  components, err := models.AllComponentsByApp(a.ID)
  if err != nil {
    lumber.Error("component:Clean:models.AllComponentsByApp(%s): %s", a.ID, err.Error())
    return fmt.Errorf("failed to fetch app component collection: %s", err.Error())
  }
  
  // iterate through the components and clean them
  for _, component := range components {
    if err := cleanComponent(a, component); err != nil {
      return fmt.Errorf("failed to clean component: %s", err.Error())
    }
  }
  
  return nil
}

// cleanComponent will clean a component if it was left in a bad state
func cleanComponent(a *models.App, component *models.Component) error {

  // short-circuit if the component is not dirty
  if !isComponentDirty(component) {
    return nil
  }
  
  if err := Destroy(a, component); err != nil {
    return fmt.Errorf("failed to remove component: %s", err.Error())
  }

	return nil
}

// isComponentDirty returns true if the container is removed or in a bad state
func isComponentDirty(component *models.Component) bool {
  // short-circuit if this service never made it to active
  if component.State != "active" {
    return true
  }

  // let's see if the container exists
  _, err := docker.GetContainer(component.ID)
  if err != nil {
    return true
  }
  
  return false
}
