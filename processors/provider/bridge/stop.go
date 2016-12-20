package bridge

import (
	"fmt"
	
	"github.com/jcelliott/lumber"

	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/display"
)

func Stop() error {

	if util.IsPrivileged() {

		// start service
		return stopService()

	} else {

		// escalate
		return reExecPrivilageStop()

	}

	return nil
}

// reExecPrivilageStart re-execs the current process with a privileged user
func reExecPrivilageStop() error {
	display.PauseTask()
	defer display.ResumeTask()

	display.PrintRequiresPrivilege("to stop the vpn")

	cmd := fmt.Sprintf("%s env bridge stop", config.NanoboxPath())

	if err := util.PrivilegeExec(cmd); err != nil {
		lumber.Error("bridge:reExecPrivilageStart:util.PrivilegeExec(%s): %s", cmd, err)
		return err
	}

	return nil
}
