package component

import (
  "fmt"
  
  "github.com/jcelliott/lumber"
  
  "github.com/nanobox-io/nanobox/models"
  "github.com/nanobox-io/nanobox-boxfile"
)

// Sync syncronizes an app's components with the boxfile config
func Sync(e *models.Env, a *models.App) error {
  
  // purge delta components
  if err := purgeDeltaComponents(e, a); err != nil {
    return fmt.Errorf("failed to purge delta components: %s", err.Error())
  }
  
  // provision components
  if err := provisionComponents(e, a); err != nil {
    return fmt.Errorf("failed to provision components: %s", err.Error())
  }
  
  // update deployed boxfile
  a.DeployedBoxfile = e.BuiltBoxfile
  if err := a.Save(); err != nil {
    lumber.Error("component:Sync:models.App.Save(): %s", err.Error())
    return fmt.Errorf("failed to update deployed boxfile on app: %s", err.Error())
  }
  
  return nil
}

// purgeDeltaComponents purges components that have changed in the boxfile
func purgeDeltaComponents(e *models.Env, a *models.App) error {
  // parse the boxfiles
  builtBoxfile    := boxfile.New([]byte(e.BuiltBoxfile))
  deployedBoxfile := boxfile.New([]byte(a.DeployedBoxfile))
  
  components, err := models.AllComponentsByApp(a.ID)
  if err != nil {
    lumber.Error("component:purgeDeltaComponents:models.AllComponentsByApp(%s): %s", a.ID, err.Error())
    return fmt.Errorf("failed to load app components: %s", err.Error())
  }
  
  for _, component := range components {
    
    // ignore platform services
    if isPlatformUID(component.Name) {
      continue
    }
    
    // fetch the data nodes
    newNode := builtBoxfile.Node(component.Name)
    oldNode := deployedBoxfile.Node(component.Name)
    
    // skip if the new node is valid and they are the same
    if newNode.Valid && newNode.Equal(oldNode) {
      continue
    }
    
    // destroy the component
    if err := Destroy(a, component); err != nil {
      return fmt.Errorf("failed to destroy component: %s", err.Error())
    }
  }
  
  return nil
}

// provisionComponents will provision components from the boxfile
func provisionComponents(e *models.Env, a *models.App) error {
  // parse the boxfile
  builtBoxfile := boxfile.New([]byte(e.BuiltBoxfile))
  
  // grab all of the data nodes
  dataServices := builtBoxfile.Nodes("data")
  
  for _, name := range dataServices {
    // check to see if this component is already active
    comp, _ := models.FindComponentBySlug(a.ID, name)
    if comp.State == "active" {
      continue
    }

    // fetch the image
    image := builtBoxfile.Node(name).StringValue("image")

    // setup
    if err := Setup(a, name, name, image); err != nil {
      return fmt.Errorf("failed to setup component (%s): %s", name, err.Error())
    }
    
    // load the component
    c, err := models.FindComponentBySlug(a.ID, name)
    if err != nil {
      lumber.Error("component:provisionComponents:models.FindComponentBySlug(%s, %s): %s", a.ID, name, err.Error())
      return fmt.Errorf("failed to load the component: %s", err.Error())
    }
    
    // configure
    if err := Configure(a, c); err != nil {
      return fmt.Errorf("failed to configure component: %s", err.Error())
    }
  }
  
  return nil
}

// isPlatform will return true if the uid matches a platform service
func isPlatformUID(uid string) bool {
	return uid == "portal" || uid == "hoarder" || uid == "mist" || uid == "logvac"
}
