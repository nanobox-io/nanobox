package component

import (
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
		return util.ErrorAppend(err, "failed to generate component data")
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
			return util.ErrorAppend(err, "unable to retrieve component image")
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
			return util.ErrorAppend(err, "failed to pull docker image (%s): %s", componentModel.Image, err.Error())
		}
		display.StopTask()
	}
	// reserve IPs
	if err := reserveIP(appModel, componentModel); err != nil {
		return util.ErrorAppend(err, "failed to reserve IPs for component")
	}

	// start the container
	display.StartTask("Starting docker container")
	config := container_generator.ComponentConfig(componentModel)

	// remove any container that may have been created with this name befor
	// this can happen if the process is killed after the 
	// container was created but before our db model was saved
	docker.ContainerRemove(config.Name)
	
	container, err := docker.CreateContainer(config)
	if err != nil {
		lumber.Error("component:Setup:docker.CreateContainer(%+v): %s", config, err.Error())
		display.ErrorTask()
		return util.ErrorAppend(err, "failed to start docker container")
	}
	display.StopTask()

	// persist the container ID
	componentModel.ID = container.ID
	if err := componentModel.Save(); err != nil {
		lumber.Error("component:Setup:models.Component.Save()")
		return util.ErrorAppend(err, "failed to persist container ID")
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
		return util.ErrorAppend(err, "failed to set component state")
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
				lumber.Error("component.reserveIPs:dhcp.ReserveLocal()")
				return util.ErrorAppend(err, "failed to reserve local IP address")
			}

			componentModel.IP = localIP.String()
		}
	}

	if err := componentModel.Save(); err != nil {
		display.StopTask()
		lumber.Error("component.reserveIPs:models.Component.Save()")
		return util.ErrorAppend(err, "failed to persist component IPs")
	}

	return nil
}

// planComponent gathers information about the componenent
func planComponent(appModel *models.App, componentModel *models.Component) error {
	display.StartTask("Gathering requirements")
	defer display.StopTask()

	planOutput, err := hookit.DebugExec(componentModel.ID, "plan", hook_generator.PlanPayload(componentModel), "info")
	if err != nil {
		return util.ErrorAppend(err, "failed to run plan hook")
	}

	// generate the component plan
	if err := componentModel.GeneratePlan(planOutput); err != nil {
		lumber.Error("component:Setup:models.Component.GeneratePlan(%s): %s", planOutput, err.Error())
		return util.ErrorAppend(err, "failed to generate the component plan")
	}

	// generate environment variables
	if err := componentModel.GenerateEvars(appModel); err != nil {
		lumber.Error("component:Setup:models.Component.GenerateEvars(%+v): %s", appModel, err.Error())
		return util.ErrorAppend(err, "failed to generate the component evars")
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
		return util.ErrorAppend(err, "failed to run update hook")
	}

	// run the configure hook
	if _, err := hookit.DebugExec(componentModel.ID, "configure", hook_generator.ConfigurePayload(appModel, componentModel), "info"); err != nil {
		return util.ErrorAppend(err, "failed to run configure hook")
	}

	// run the start hook
	if _, err := hookit.DebugExec(componentModel.ID, "start", hook_generator.UpdatePayload(componentModel), "info"); err != nil {
		return util.ErrorAppend(err, "failed to run start hook")
	}

	return nil
}
