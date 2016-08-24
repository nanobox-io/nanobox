package app

import (
  "fmt"
  
  "github.com/jcelliott/lumber"
  
  "github.com/nanobox-io/nanobox/models"
  "github.com/nanobox-io/nanobox/processors/component"
  "github.com/nanobox-io/nanobox/processors/platform"
  "github.com/nanobox-io/nanobox/util/dhcp"
)

// Setup sets up the app on the provider and in the database
func Setup(e *models.Env, a *models.App, name string) error {
	locker.LocalLock()
	defer locker.LocalUnlock()  
  
  // short-circuit if this app is already active
  if a.State == "active" {
    return nil
  }
  
  // generate the app data
  if err := a.Generate(e, name); err != nil {
    lumber.Error("app:Setup:models.App:Generate(): %s", err.Error())
    return fmt.Errorf("failed to generate app data: %s", err.Error())
  }
  
  // reserve IPs
  if err := reserveIPs(a); err != nil {
    return fmt.Errorf("failed to reserve app IPs: %s", err.Error())
  }

  // clean crufty components
  if err := component.Clean(a); err != nil {
    return fmt.Errorf("failed to clean crufty components: %s", err.Error())
  }
  
  // setup the platform services
  if err := platform.Setup(a); err != nil {
    return fmt.Errorf("failed to setup platform services: %s", err.Error())
  }

  // set app state to active
  a.State = "active"
  if err := a.Save(); err != nil {
    lumber.Error("app:Setup:models:App:Save(): %s", err.Error())
    return fmt.Errorf("failed to persist app state: %s", err.Error())
  }

  return nil
}

// reserIPs reserves app-level ip addresses
func reserveIPs(a *models.App) error {
  // reserve a dev ip
  envIP, err := dhcp.ReserveGlobal()
  if err != nil {
    lumber.Error("app:reserveIPs:dhcp.ReserveGlobal(): %s", err.Error())
    return fmt.Errorf("failed to reserve an env IP: %s", err.Error())
  }

  // reserve a logvac ip
  logvacIP, err := dhcp.ReserveLocal()
  if err != nil {
    lumber.Error("app:reserveIPs:dhcp.ReserveLocal(): %s", err.Error())
    return fmt.Errorf("failed to reserve a logvac IP: %s", err.Error())
  }

  // reserve a mist ip
  mistIP, err := dhcp.ReserveLocal()
  if err != nil {
    lumber.Error("app:reserveIPs:dhcp.ReserveLocal(): %s", err.Error())
    return fmt.Errorf("failed to reserve a mist IP: %s", err.Error())
  }
  
  // now assign the IPs onto the app model
  a.GlobalIPs["env"] = envIP.String()
  
  a.LocalIPs["logvac"] = logvacIP.String()
  a.LocalIPs["mist"]   = mistIP.String()
  
  // save the app
  if err := a.Save(); err != nil {
    lumber.Error("app:reserveIPs:models:App:Save(): %s", err.Error())
    return fmt.Errorf("failed to persist IPs to the db: %s", err.Error())
  }
  
  return nil
}
