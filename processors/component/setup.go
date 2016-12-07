package component

import (
	"fmt"
	"time"

	"github.com/jcelliott/lumber"
	"github.com/nanobox-io/golang-docker-client"

	container_generator "github.com/nanobox-io/nanobox/generators/containers"
	hook_generator "github.com/nanobox-io/nanobox/generators/hooks/component"
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/dhcp"
	"github.com/nanobox-io/nanobox/util/display"
	"github.com/nanobox-io/nanobox/util/hookit"
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
		image, err := componentImage(componentModel)
		if err != nil {
			lumber.Error("component:Setup:boxfile.ComponentImage(%+v): %s", componentModel, err.Error())
			return fmt.Errorf("unable to retrieve component image: %s", err.Error())
		}
		componentModel.Image = image
	}

	// generate a docker percent display
	dockerPercent := &display.DockerPercentDisplay{
		Output: display.NewStreamer("info"),
		// Prefix: componentModel.Image,
	}

	if !docker.ImageExists(componentModel.Image) {

		// pull the component image
		display.StartTask("Pulling %s image", componentModel.Image)
		imagePull := func() error {
			_, err := docker.ImagePull(componentModel.Image, dockerPercent)
			return err
		}
		if err := util.Retry(imagePull, 5, time.Second); err != nil {
			lumber.Error("component:Setup:docker.ImagePull(%s, nil): %s", componentModel.Image, err.Error())
			// remove the component because it doesnt need to be cleaned up at this point
			componentModel.Delete()
			display.ErrorTask()
			return fmt.Errorf("failed to pull docker image (%s): %s", componentModel.Image, err.Error())
		}
		display.StopTask()
	}
	// reserve IPs
	if err := reserveIP(appModel, componentModel); err != nil {
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

// reserveIP reserves IP addresses for this component
func reserveIP(appModel *models.App, componentModel *models.Component) error {
	display.StartTask("Reserve IP")
	defer display.StopTask()

	// dont reserve a new one if we already have this one
	if componentModel.IPAddr() == "" {
		// first let's see if our local IP was reserved during app creation
		if componentModel.Name == "portal" {
			componentModel.IP = appModel.LocalIPs["env"]
		} else if appModel.LocalIPs[componentModel.Name] != "" {

			// assign the localIP from the pre-generated app cache
			componentModel.IP = appModel.LocalIPs[componentModel.Name]
		} else {

			localIP, err := dhcp.ReserveLocal()
			if err != nil {
				display.StopTask()
				lumber.Error("component.reserveIPs:dhcp.ReserveLocal(): %s", err.Error())
				return fmt.Errorf("failed to reserve local IP address: %s", err.Error())
			}

			componentModel.IP = localIP.String()
		}
	}

	if err := componentModel.Save(); err != nil {
		display.StopTask()
		lumber.Error("component.reserveIPs:models.Component.Save(): %s", err.Error())
		return fmt.Errorf("failed to persist component IPs: %s", err.Error())
	}

	return nil
}

// planComponent gathers information about the componenent
func planComponent(appModel *models.App, componentModel *models.Component) error {
	display.StartTask("Gathering requirements")
	defer display.StopTask()

	planOutput, err := hookit.DebugExec(componentModel.ID, "plan", hook_generator.PlanPayload(componentModel), "info")
	if err != nil {
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
	if _, err := hookit.DebugExec(componentModel.ID, "update", hook_generator.UpdatePayload(componentModel), "debug"); err != nil {
		display.ErrorTask()
		return fmt.Errorf("failed to run update hook: %s", err.Error())
	}

	// run the configure hook
	if _, err := hookit.DebugExec(componentModel.ID, "configure", hook_generator.ConfigurePayload(appModel, componentModel), "info"); err != nil {
		return fmt.Errorf("failed to run configure hook: %s", err.Error())
	}

	// run the start hook
	if _, err := hookit.DebugExec(componentModel.ID, "start", hook_generator.UpdatePayload(componentModel), "info"); err != nil {
		return fmt.Errorf("failed to run start hook: %s", err.Error())
	}

	return nil
}
