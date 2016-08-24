package app

import (
  "fmt"
  "net"
  
  "github.com/jcelliott/lumber"
  "github.com/nanobox-io/golang-docker-client"
  
  "github.com/nanobox-io/nanobox/models"
  "github.com/nanobox-io/nanobox/processors/app/dns"
  "github.com/nanobox-io/nanobox/processors/component"
  "github.com/nanobox-io/nanobox/processors/provider"
  "github.com/nanobox-io/nanobox/util/dhcp"
)

// Destroy removes the app from the provider and the database
func Destroy(a *models.App) error {
  locker.LocalLock()
  defer locker.LocalUnlock()
  
  // short-circuit if this app isn't created
  if a.IsNew() {
    return nil
  }
  
  // initialize docker for the provider
  if err := provider.Init(); err != nil {
    return fmt.Errorf("failed to initialize docker environment: %s", err.Error())
  }
  
  // remove the dev container if there is one
  docker.ContainerRemove(fmt.Sprintf("nanobox_%s", a.ID))
  
  // destroy the associated components
  if err := destroyComponents(s); err != nil {
    return fmt.Errorf("failed to destroy components: %s", err.Error())
  }
  
  // release IPs
  if err := releaseIPs(a); err != nil {
    return fmt.Errorf("failed to release IPs: %s", err.Error())
  }
  
  // remove dns entries for this app
  if err := dns.RemoveAll(a); err != nil {
    return fmt.Errorf("failed to clean dns entries: %s", err.Error())
  }
  
  // destroy the app model
  if err := a.Delete(); err != nil {
    lumber.Error("app:Destroy:models.App.Destroy(): %s", err.Error())
    return fmt.Errorf("failed to delete app model: %s", err.Error())
  }
  
  return nil
}

// destroyComponents destroys all the components of this app
func destroyComponents(a *models.App) error {
  components, err := models.AllComponentsByApp(a.ID)
  if err != nil {
    lumber.Error("app:destroyComponents:models.AllComponentsByApp(%s) %s", a.ID, err.Error())
    return fmt.Errorf("unable to retrieve components: %s", err.Error())
  }
  
  for _, c := range components {
    if err := component.Destroy(a, c); err != nil {
      return fmt.Errorf("failed to destroy app component: %s", err.Error())
    }
  }
  
  return nil
}

// releaseIPs releases the app-level ip addresses
func releaseIPs(a *models.App) error {
  // release all of the external IPs
  for _, ip := range a.GlobalIPs {
    // release the IP
    if err := dhcp.ReturnIP(net.ParseIP(ip)); err != nil {
      lumber.Error("app:Destroy:releaseIPs:dhcp.ReturnIP(%s): %s", ip, err.Error())
      return fmt.Errorf("failed to release IP: %s", err.Error())
    }
  }

  // release all of the local IPs
  for _, ip := range a.LocalIPs {
    // release the IP
    if err := dhcp.ReturnIP(net.ParseIP(ip)); err != nil {
      lumber.Error("app:Destroy:releaseIPs:dhcp.ReturnIP(%s): %s", ip, err.Error())
      return fmt.Errorf("failed to release IP: %s", err.Error())
    }
  }

  return nil
}
