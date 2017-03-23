package server

import (
	"fmt"

	"github.com/jcelliott/lumber"

	"github.com/nanobox-io/nanobox/util"	
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/display"	
	"github.com/nanobox-io/nanobox/util/service"	
)

func Stop() error {
	// run as admin
	if !util.IsPrivileged() {
		return reExecPrivilageStop()
	}

	return service.Stop("nanobox-server")
}

// reExecPrivilageStop re-execs the current process with a privileged user
func reExecPrivilageStop() error {
	display.PauseTask()
	defer display.ResumeTask()

	display.PrintRequiresPrivilege("to stop the server")

	cmd := fmt.Sprintf("%s env server stop", config.NanoboxPath())

	if err := util.PrivilegeExec(cmd); err != nil {
		lumber.Error("server:reExecPrivilageStop:util.PrivilegeExec(%s): %s", cmd, err)
		return err
	}

	return nil
}
