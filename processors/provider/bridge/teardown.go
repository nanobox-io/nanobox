package bridge

import (
	"fmt"
	"os"
	"runtime"

	"github.com/jcelliott/lumber"
	"github.com/nanobox-io/golang-docker-client"

	container_generator "github.com/nanobox-io/nanobox/generators/containers"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/display"
	"github.com/nanobox-io/nanobox/util/config"
)

func Teardown() error {

	if util.IsPrivileged() {
		// remove bridge client
		if err := Stop(); err != nil {
			return err
		}

		// remove bridge config
		if err := removeConfig(); err != nil {
			return err
		}

	} else {
		// remove component
		if err := removeComponent(); err != nil {
			return err
		}
		return reExecTeardown()
	}

	return nil
}

func removeConfig() error {
	if runtime.GOOS != "windows" {
		return os.Remove(serviceConfigFile())
	}

	return nil
}

func removeComponent() error {
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


// reExecTeardown re-execs the current process with a privileged user
func reExecTeardown() error {
	display.PauseTask()
	defer display.ResumeTask()

	display.PrintRequiresPrivilege("remove the vpn")

	cmd := fmt.Sprintf("%s env bridge teardown", config.NanoboxPath())

	if err := util.PrivilegeExec(cmd); err != nil {
		lumber.Error("bridge:reExecTeardown:util.PrivilegeExec(%s): %s", cmd, err)
		return err
	}

	return nil
}
