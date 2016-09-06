package component

import (
	"fmt"

	"github.com/jcelliott/lumber"
	"github.com/nanobox-io/golang-docker-client"

	container_generator "github.com/nanobox-io/nanobox/generators/containers"
	hook_generator "github.com/nanobox-io/nanobox/generators/hooks/component"
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util/boxfile"
	"github.com/nanobox-io/nanobox/util/dhcp"
	"github.com/nanobox-io/nanobox/util/display"
	"github.com/nanobox-io/nanobox/util/hookit"
	"github.com/nanobox-io/nanobox/util/provider"
)

// Setup sets up the component container and model data
func Setup(appModel *models.App, componentModel *models.Component) error {

	// generate the missing component data
	if err := componentModel.Generate(appModel, "data"); err != nil {
		lumber.Error("component:Setup:models.Component:Generate(%s, data): %s", appModel.ID, componentModel.Name, err.Error())
		return fmt.Errorf("failed to generate component data: %s", err.Error())
	}

	// short-circuit if this component is already setup
	if componentModel.State != "initialized" {
		return nil
	}

	display.OpenContext(componentModel.Label)
	defer display.CloseContext()

	// if the image was not provided
	if componentModel.Image == "" {
		// extract the image from the boxfile node
		image, err := boxfile.ComponentImage(componentModel)
		if err != nil {
			lumber.Error("component:Setup:boxfile.ComponentImage(%+v): %s", componentModel, err.Error())
			return fmt.Errorf("unable to retrieve component image: %s", err.Error())
		}
		componentModel.Image = image
	}

	// generate a docker percent display
	dockerPercent := &display.DockerPercentDisplay{
		Output: display.NewStreamer("info"),
		Prefix: componentModel.Image,
	}

	// pull the component image
	display.StartTask("Pulling %s image", componentModel.Image)
	if _, err := docker.ImagePull(componentModel.Image, dockerPercent); err != nil {
		lumber.Error("component:Setup:docker.ImagePull(%s, nil): %s", componentModel.Image, err.Error())
		display.ErrorTask()
		return fmt.Errorf("failed to pull docker image (%s): %s", componentModel.Image, err.Error())
	}
	display.StopTask()

	// reserve IPs
	if err := reserveIPs(appModel, componentModel); err != nil {
		return fmt.Errorf("failed to reserve IPs for component: %s", err.Error())
	}

	// start the container
	display.StartTask("Starting docker container")
	config := container_generator.ComponentConfig(componentModel)
	container, err := docker.CreateContainer(config)
	if err != nil {
		lumber.Error("component:Setup:docker.CreateContainer(%+v): %s", config, err.Error())
		display.ErrorTask()
		return fmt.Errorf("failed to start docker container: %s", err.Error())
	}
	display.StopTask()

	// persist the container ID
	componentModel.ID = container.ID
	if err := componentModel.Save(); err != nil {
		lumber.Error("component:Setup:models.Component.Save(): %s", err.Error())
		return fmt.Errorf("failed to persist container ID: %s", err.Error())
	}

	// attach container to the host network
	if err := attachNetwork(componentModel); err != nil {
		return fmt.Errorf("failed to attach container to host network: %s", err.Error())
	}

	// plan the component
	if err := planComponent(appModel, componentModel); err != nil {
		return err
	}

	if err := configureComponent(appModel, componentModel); err != nil {
		return err
	}

	// set state as active
	componentModel.State = "active"
	if err := componentModel.Save(); err != nil {
		lumber.Error("component:Setup:models.Component.Save()", err.Error())
		return fmt.Errorf("failed to set component state: %s", err.Error())
	}

	return nil
}

// reserveIPs reserves IP addresses for this component
func reserveIPs(appModel *models.App, componentModel *models.Component) error {
	display.StartTask("Reserve IPs")
	defer display.StopTask()

	// dont reserve a new one if we already have this one
	if componentModel.InternalIP == "" {
		// first let's see if our local IP was reserved during app creation
		if appModel.LocalIPs[componentModel.Name] != "" {

			// assign the localIP from the pre-generated app cache
			componentModel.InternalIP = appModel.LocalIPs[componentModel.Name]
		} else {

			localIP, err := dhcp.ReserveLocal()
			if err != nil {
				display.StopTask()
				lumber.Error("component.reserveIPs:dhcp.ReserveLocal(): %s", err.Error())
				return fmt.Errorf("failed to reserve local IP address: %s", err.Error())
			}

			componentModel.InternalIP = localIP.String()
		}
	}

	// dont reserve a new global ip if i already have one
	if componentModel.ExternalIP == "" {
		// only if this service is portal, we need to use the preview IP
		// in a dev environment there will be no portal installed
		// so the env ip should be available
		// in dev the env ip is used for the dev container
		if componentModel.Name == "portal" {
			// portal's global ip is the preview ip
			componentModel.ExternalIP = appModel.GlobalIPs["env"]
		} else {

			globalIP, err := dhcp.ReserveGlobal()
			if err != nil {
				display.StopTask()
				lumber.Error("component.reserveIPs:dhcp.ReserveGlobal(): %s", err.Error())
				return fmt.Errorf("failed to reserve global IP address: %s", err.Error())
			}

			componentModel.ExternalIP = globalIP.String()
		}
	}

	if err := componentModel.Save(); err != nil {
		display.StopTask()
		lumber.Error("component.reserveIPs:models.Component.Save(): %s", err.Error())
		return fmt.Errorf("failed to persist component IPs: %s", err.Error())
	}

	return nil
}

// attachNetwork attaches the component to the host network
func attachNetwork(componentModel *models.Component) error {
	display.StartTask("Attaching network")
	defer display.StopTask()

	// add the IP to the provider
	if err := provider.AddIP(componentModel.ExternalIP); err != nil {
		lumber.Error("component:Setup:attachNetwork:provider.AddIP(%s): %s", componentModel.ExternalIP, err.Error())
		display.ErrorTask()
		return fmt.Errorf("failed to add IP to provider: %s", err.Error())
	}

	// nat traffic from the external IP to the internal
	if err := provider.AddNat(componentModel.ExternalIP, componentModel.InternalIP); err != nil {
		lumber.Error("component:Setup:attachNetwork:provider.AddNat(%s, %s): %s", componentModel.ExternalIP, componentModel.InternalIP, err.Error())
		display.ErrorTask()
		return fmt.Errorf("failed to nat IP on provider: %s", err.Error())
	}

	return nil
}

// planComponent gathers information about the componenent
func planComponent(appModel *models.App, componentModel *models.Component) error {
	display.StartTask("Gathering requirements")
	defer display.StopTask()

	planOutput, err := hookit.RunPlanHook(componentModel.ID, hook_generator.PlanPayload(componentModel))
	if err != nil {
		display.ErrorTask()
		return fmt.Errorf("failed to run plan hook: %s", err.Error())
	}

	// generate the component plan
	if err := componentModel.GeneratePlan(planOutput); err != nil {
		lumber.Error("component:Setup:models.Component.GeneratePlan(%s): %s", planOutput, err.Error())
		return fmt.Errorf("failed to generate the component plan: %s", err.Error())
	}

	// generate environment variables
	if err := componentModel.GenerateEvars(appModel); err != nil {
		lumber.Error("component:Setup:models.Component.GenerateEvars(%+v): %s", appModel, err.Error())
		return fmt.Errorf("failed to generate the component evars: %s", err.Error())
	}

	return nil
}

// configureComponent configures the component
func configureComponent(appModel *models.App, componentModel *models.Component) error {
	display.StartTask("Configuring services")
	defer display.StopTask()

	// run the update hook
	if _, err := hookit.RunUpdateHook(componentModel.ID, hook_generator.UpdatePayload(componentModel)); err != nil {
		display.ErrorTask()
		return fmt.Errorf("failed to run update hook: %s", err.Error())
	}

	// run the configure hook
	if _, err := hookit.RunConfigureHook(componentModel.ID, hook_generator.ConfigurePayload(appModel, componentModel)); err != nil {
		display.ErrorTask()
		return fmt.Errorf("failed to run configure hook: %s", err.Error())
	}

	// run the start hook
	if _, err := hookit.RunStartHook(componentModel.ID, hook_generator.UpdatePayload(componentModel)); err != nil {
		display.ErrorTask()
		return fmt.Errorf("failed to run start hook: %s", err.Error())
	}

	return nil
}
