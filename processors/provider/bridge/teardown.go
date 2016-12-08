package bridge

import (
	"fmt"
	
	"github.com/jcelliott/lumber"
	"github.com/nanobox-io/golang-docker-client"

	container_generator "github.com/nanobox-io/nanobox/generators/containers"
	"github.com/nanobox-io/nanobox/util/locker"

)


func Teardown() error {
	// remove bridge client
	if err := removeBridge(); err != nil {
		return err
	}

	// remove bridge config
	if err := removeConfig(); err != nil {
		return err
	}

	// remove component
	if err := removeComponent(); err != nil {
		return err
	}
	
	return nil	
}

func removeBridge() error {
	
	return nil
}


func removeConfig() error {
	
	return nil
}

func removeComponent() error {
	locker.LocalLock()
	defer locker.LocalUnlock()

	// grab the container info
	container, err := docker.GetContainer(container_generator.BridgeName())
	if err != nil {
		// if we cant get the container it may have been removed by someone else
		// just return here
		return nil
	}

	// remove the container
	if err := docker.ContainerRemove(container.ID); err != nil {
		lumber.Error("provider:bridge:teardown:docker.ContainerRemove(%s): %s", container.ID, err)
		return fmt.Errorf("failed to remove bridge container: %s", err.Error())
	}

	return nil
}
