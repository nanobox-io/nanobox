package component

import (
  "fmt"
  
  "github.com/jcelliott/lumber"
  "github.com/nanobox-io/golang-docker-client"
  
  "github.com/nanobox-io/nanobox/generators/containers"
  "github.com/nanobox-io/nanobox/models"
  "github.com/nanobox-io/nanobox/util/boxfile"
  "github.com/nanobox-io/nanobox/util/display"
  "github.com/nanobox-io/nanobox/util/dhcp"
  "github.com/nanobox-io/nanobox/util/provider"
)

// Setup sets up the component container and model data
func Setup(a *models.App, name, label, image string) error {
  // fetch the component or get a new one
  c, _ := models.FindComponentBySlug(a.ID, name)
  
  // generate the component data
  if err := c.Generate(a, name, label, image, "data"); err != nil {
    lumber.Error("component:Setup:models.Component:Generate(%s, %s, data): %s", a.ID, name, err.Error())
    return fmt.Errorf("failed to generate component data: %s", err.Error())
  }
  
  // short-circuit if this component is already setup
  if c.State != "initialized" {
    return nil
  }
  
  // extract the image from the boxfile node
  image, err := boxfile.ComponentImage(c)
  if err != nil {
    lumber.Error("component:Setup:boxfile.ComponentImage(%+v): %s", c, err.Error())
    return fmt.Errorf("unable to retrieve component image: %s", err.Error())
  }
  
  // generate a docker percent display
  dockerPercent := &display.DockerPercentDisplay{
    Output: display.NewStreamer("info"), 
    Prefix: image,
  }
  
  // pull the component image
  if _, err := docker.ImagePull(image, dockerPercent); err != nil {
    lumber.Error("component:Setup:docker.ImagePull(%s, nil): %s", image, err.Error())
    return fmt.Errorf("failed to pull docker image (%s): %s", image, err.Error())
  }
  
  // reserve IPs
  if err := reserveIPs(a, c); err != nil {
    return fmt.Errorf("failed to reserve IPs for component: %s", err.Error())
  }
  
  // start the container
  config := containers.ComponentConfig(c, image, c.InternalIP)
  container, err := docker.CreateContainer(config)
  if err != nil {
    lumber.Error("component:Setup:docker.CreateContainer(%+v): %s", config, err.Error())
    return fmt.Errorf("failed to start docker container: %s", err.Error())
  }
  
  // persist the container ID
  c.ID = container.ID
  if err := c.Save(); err != nil {
    lumber.Error("component:Setup:models.Component.Save(): %s", err.Error())
    return fmt.Errorf("failed to persist container ID: %s", err.Error())
  }
  
  // attach container to the host network
  if err := attachNetwork(c); err != nil {
    return fmt.Errorf("failed to attach container to host network: %s", err.Error())
  }
  
  // plan the component
  planOutput, err := RunPlanHook(c)
  if err != nil {
    return fmt.Errorf("failed to run plan hook: %s", err.Error())
  }
  
  // generate the component plan
  if err := c.GeneratePlan(planOutput); err != nil {
    lumber.Error("component:Setup:models.Component.GeneratePlan(%s): %s", planOutput, err.Error())
    return fmt.Errorf("failed to generate the component plan: %s", err.Error())
  }
  
  // generate environment variables
  if err := c.GenerateEvars(a); err != nil {
    lumber.Error("component:Setup:models.Component.GenerateEvars(%+v): %s", a, err.Error())
    return fmt.Errorf("failed to generate the component evars: %s", err.Error())
  }
  
  // update state
  c.State = "planned"
  if err := c.Save(); err != nil {
    lumber.Error("component:Setup:models.Component.Save(): %s", err.Error())
    return fmt.Errorf("failed to persist component state: %s", err.Error())
  }
  
  return nil
}

// reserveIPs reserves IP addresses for this component
func reserveIPs(a *models.App, c *models.Component) error {
  // dont reserve a new one if we already have this one
  if c.InternalIP == "" {
    // first let's see if our local IP was reserved during app creation
    if a.LocalIPs[c.Name] != "" {

      // assign the localIP from the pre-generated app cache
      c.InternalIP = a.LocalIPs[c.Name]
    } else {

      localIP, err := dhcp.ReserveLocal()
      if err != nil {
        lumber.Error("component.reserveIPs:dhcp.ReserveLocal(): %s", err.Error())
        return fmt.Errorf("failed to reserve local IP address: %s", err.Error())
      }

      c.InternalIP = localIP.String()
    }
  }

  // dont reserve a new global ip if i already have one
  if c.ExternalIP == "" {
    // only if this service is portal, we need to use the preview IP
    // in a dev environment there will be no portal installed
    // so the env ip should be available
    // in dev the env ip is used for the dev container
    if c.Name == "portal" {
      // portal's global ip is the preview ip
      c.ExternalIP = a.GlobalIPs["env"]
    } else {

      globalIP, err := dhcp.ReserveGlobal()
      if err != nil {
        lumber.Error("component.reserveIPs:dhcp.ReserveGlobal(): %s", err.Error())
        return fmt.Errorf("failed to reserve global IP address: %s", err.Error())
      }

      c.ExternalIP = globalIP.String()
    }
  }
  
  if err := c.Save(); err != nil {
    lumber.Error("component.reserveIPs:models.Component.Save(): %s", err.Error())
    return fmt.Errorf("failed to persist component IPs: %s", err.Error())
  }
  
  return nil
}

// attachNetwork attaches the component to the host network
func attachNetwork(c *models.Component) error {
  // add the IP to the provider
  if err := provider.AddIP(c.ExternalIP); err != nil {
    lumber.Error("component:Setup:attachNetwork:provider.AddIP(%s): %s", c.ExternalIP, err.Error())
    return fmt.Errorf("failed to add IP to provider: %s", err.Error())
  }
  
  // nat traffic from the external IP to the internal
  if err := provider.AddNat(c.ExternalIP, c.InternalIP); err != nil {
    lumber.Error("component:Setup:attachNetwork:provider.AddNat(%s, %s): %s", c.ExternalIP, c.InternalIP, err.Error())
    return fmt.Errorf("failed to nat IP on provider: %s", err.Error())
  }
  
  return nil
}
