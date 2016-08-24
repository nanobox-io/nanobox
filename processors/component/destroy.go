package component

import (
  "fmt"
  "net"
  
  "github.com/jcelliott/lumber"
  "github.com/nanobox-io/golang-docker-client"  
  
  "github.com/nanobox-io/nanobox/models"
  "github.com/nanobox-io/nanobox/util/provider"
  "github.com/nanobox-io/nanobox/util/dhcp"
)

// Destroy destroys a component from the provider and database
func Destroy(a *models.App, c *models.Component) error {
  
  // remove the docker container
  if err := docker.ContainerRemove(c.ID); err != nil {
    lumber.Error("component:Destroy:docker.ContainerRemove(%s): %s", c.ID, err.Error())
    return fmt.Errorf("failed to remove docker container: %s", err.Error())
  }
  
  // detach from the host network
  if err := detachNetwork(a, c); err != nil {
    return fmt.Errorf("failed to detach container from the host network: %s", err.Error())
  }
  
  // purge evars
  if err := c.PurgeEvars(a); err != nil {
    lumber.Error("component:Destroy:models.Component.PurgeEvars(%+v): %s", a, err.Error())
    return fmt.Errorf("failed to purge component evars from app: %s", err.Error())
  }
  
  // destroy the data model
  if err := c.Delete(); err != nil {
    lumber.Error("component:Destroy:models.Component.Delete(): %s", err.Error())
    return fmt.Errorf("failed to destroy component model: %s", err.Error())
  }
  
  return nil
}

// detachNetwork detaches the network from the host
func detachNetwork(a *models.App, c *models.Component) error {
  
  // remove NAT
  if err := provider.RemoveNat(c.ExternalIP, c.InternalIP); err != nil {
    lumber.Error("component:detachNetwork:provider.RemoveNat(%s, %s): %s", c.ExternalIP, c.InternalIP, err.Error())
    return fmt.Errorf("failed to remove NAT from provider: %s", err.Error())
  }
  
  // remove IP
  if err := provider.RemoveIP(c.ExternalIP); err != nil {
    lumber.Error("component:detachNetwork:provider.RemoveIP(%s): %s", c.ExternalIP, err.Error())
    return fmt.Errorf("failed to remove IP from provider: %s", err.Error())
  }
  
  // return the external IP
  // don't return the external IP if this is portal
  if c.Name != "portal" && a.GlobalIPs[c.Name] == "" {
    ip := net.ParseIP(c.ExternalIP)
    if err := dhcp.ReturnIP(ip); err != nil {
      lumber.Error("component:detachNetwork:dhcp.ReturnIP(%s): %s", ip, err.Error())
      return fmt.Errorf("failed to release IP back to pool: %s", err.Error())
    }
  }
  
  // return the internal IP
  // don't return the internal IP if it's an app-level cache
  if a.LocalIPs[c.Name] == "" {
    ip := net.ParseIP(c.InternalIP)
    if err := dhcp.ReturnIP(ip); err != nil {
      lumber.Error("component:detachNetwork:dhcp.ReturnIP(%s): %s", ip, err.Error())
      return fmt.Errorf("failed to release IP back to pool: %s", err.Error())
    }
  }
  
  return nil
}
