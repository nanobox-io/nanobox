package code

import (
	"fmt"

	"github.com/jcelliott/lumber"
	"github.com/nanobox-io/golang-docker-client"

	
	"github.com/nanobox-io/nanobox/generators/containers"
	"github.com/nanobox-io/nanobox/generators/hooks/build"
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/dhcp"
)

// Publish ...
func Publish(envModel *models.Env, WarehouseConfig WarehouseConfig) error {

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
	config := generate_container.BuildConfig(buildImage, ip.String())
	container, err := docker.CreateContainer(config)
	if err != nil {
		lumber.Error("code:Build:docker.CreateContainer(%+v): %s", config, err.Error())
		return fmt.Errorf("failed to start docker container: %s", err.Error())
	}

	// ensure we stop the container when we're done
	defer docker.ContainerRemove(container.ID)

	lumber.Prefix("code:Publish")
	defer lumber.Prefix("")

	// run user hook
	// TODO: display output from hooks
	payload, err := build.UserPayload()
	if err != nil {
		lumber.Error("code:Publish:build.UserPayload(): %s", err.Error())
		return fmt.Errorf("unable to retrieve user payload: %s", err.Error())
	}
	if _, err := util.Exec(container.ID, "user", payload, nil); err != nil {
		return runDebugSession(container.ID, err)
	}

	// TODO: i dont like this, passing too many things makes this very messy
	payload = build.PublishPayload(
		WarehouseConfig.BuildID,
		WarehouseConfig.WarehouseURL,
		WarehouseConfig.WarehouseToken,
		envModel.BuiltBoxfile,
		WarehouseConfig.PreviousBuild)
	if err != nil {
		lumber.Error("code:Publish:build.UserPayload(): %s", err.Error())
		return fmt.Errorf("unable to retrieve user payload: %s", err.Error())
	}
	if _, err := util.Exec(container.ID, "publish", payload, nil); err != nil {
		return runDebugSession(container.ID, err)
	}

	return nil
}
