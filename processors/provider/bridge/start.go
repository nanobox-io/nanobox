package bridge

import (
	"fmt"

	"github.com/jcelliott/lumber"

	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/display"
	"github.com/nanobox-io/nanobox/util/provider/bridge"
)

func Start() error {

	if util.IsPrivileged() {

		// create service
		if err := bridge.CreateService(); err != nil {
			return util.ErrorAppend(err, "failed to create service")
		}

		// start service
		if err := bridge.StartService(); err != nil {
			return util.ErrorAppend(err, "failed to start service")
		}

	} else {

		// escalate
		return reExecPrivilageStart()

	}

	return nil
}

// reExecPrivilageStart re-execs the current process with a privileged user
func reExecPrivilageStart() error {
	display.PauseTask()
	defer display.ResumeTask()

	display.PrintRequiresPrivilege("to start the vpn")

	cmd := fmt.Sprintf("%s env bridge start", config.NanoboxPath())

	if err := util.PrivilegeExec(cmd); err != nil {
		lumber.Error("bridge:reExecPrivilageStart:util.PrivilegeExec(%s): %s", cmd, err)
		return err
	}

	return nil
}
