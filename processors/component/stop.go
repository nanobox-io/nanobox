package component

import (
  "fmt"
  
  "github.com/jcelliott/lumber"
  "github.com/nanobox-io/golang-docker-client"
  
  "github.com/nanobox-io/nanobox/models"
  "github.com/nanobox-io/nanobox/util/provider"
)

// Stop stops the component's docker container
func Stop(c *models.Component) error {
  // short-circuit if the process is already stopped
  if !isComponentRunning(c.ID) {
    return nil
  }
  
  // stop the docker container
  if err := docker.ContainerStop(c.ID); err != nil {
    lumber.Error("component:Stop:docker.ContainerStop(%s): %s", c.ID, err.Error())
    return fmt.Errorf("failed to stop docker container: %s", err.Error())
  }
  
  // remove NAT
  if err := provider.RemoveNat(c.ExternalIP, c.InternalIP); err != nil {
    lumber.Error("component:Stop:provider.RemoveNat(%s, %s): %s", c.ExternalIP, c.InternalIP, err.Error())
    return fmt.Errorf("failed to remove NAT on the provider: %s", err.Error())
  }
  
  // remove the IP from the provider
  if err := provider.RemoveIP(c.ExternalIP); err != nil {
    lumber.Error("component:Stop:provider.RemoveIP(%s): %s", c.ExternalIP, err.Error())
    return fmt.Errorf("failed to remove IP from the provider: %s", err.Error())
  }
  
  // todo: set status
  
  return nil
}
