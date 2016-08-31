package code

import (
	"fmt"

	"github.com/jcelliott/lumber"
	"github.com/nanobox-io/golang-docker-client"

	container_generator "github.com/nanobox-io/nanobox/generators/containers"
	"github.com/nanobox-io/nanobox/generators/hooks/build"
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util/dhcp"
	"github.com/nanobox-io/nanobox/util/display"
	"github.com/nanobox-io/nanobox/util/hookit"
)

// Publish ...
func Publish(envModel *models.Env, WarehouseConfig WarehouseConfig) error {
	display.OpenContext("publishing code to the warehouse")
	defer display.CloseContext()

	// pull the latest build image
	buildImage, err := pullBuildImage()
	if err != nil {
		return fmt.Errorf("failed to pull the build image: %s", err.Error())
	}

	// reserve an ip
	ip, err := dhcp.ReserveLocal()
	if err != nil {
		lumber.Error("code:Publish:dhcp.ReserveLocal(): %s", err.Error())
		return err
	}
	defer dhcp.ReturnIP(ip)

	// start the container
	display.StartTask("starting publish container")
	config := container_generator.BuildConfig(buildImage, ip.String())
	container, err := docker.CreateContainer(config)
	if err != nil {
		lumber.Error("code:Build:docker.CreateContainer(%+v): %s", config, err.Error())
		display.ErrorTask()
		return fmt.Errorf("failed to start docker container: %s", err.Error())
	}
	display.StopTask()

	// ensure we stop the container when we're done
	defer docker.ContainerRemove(container.ID)

	lumber.Prefix("code:Publish")
	defer lumber.Prefix("")

	// run user hook
	// TODO: display output from hooks
	payload := build.UserPayload()
	if err != nil {
		lumber.Error("code:Publish:build.UserPayload(): %s", err.Error())
		return fmt.Errorf("unable to retrieve user payload: %s", err.Error())
	}
	if _, err := hookit.Exec(container.ID, "user", payload, "info"); err != nil {
		return runDebugSession(container.ID, err)
	}

	display.StartTask("publishing")
	buildWarehouseConfig := build.WarehouseConfig{
		BuildID: WarehouseConfig.BuildID,
		WarehouseURL: WarehouseConfig.WarehouseURL,
		WarehouseToken: WarehouseConfig.WarehouseToken,
		PreviousBuild: WarehouseConfig.PreviousBuild,
	}

	payload = build.PublishPayload(envModel, buildWarehouseConfig)
	if err != nil {
		lumber.Error("code:Publish:build.UserPayload(): %s", err.Error())
		display.ErrorTask()
		return fmt.Errorf("unable to retrieve user payload: %s", err.Error())
	}
	if _, err := hookit.Exec(container.ID, "publish", payload, "info"); err != nil {
		display.ErrorTask()
		err = fmt.Errorf("failed to run publish hook: %s", err.Error())
		return runDebugSession(container.ID, err)
	}
	display.StopTask()

	return nil
}
