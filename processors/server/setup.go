package server

import (
	"fmt"

	"github.com/jcelliott/lumber"

	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/display"
	"github.com/nanobox-io/nanobox/util/service"

)

func Setup() error {
	if service.Running("nanobox-server") {
		return nil
	}

	// run as admin
	if !util.IsPrivileged() {
		return reExecPrivilageStart()
	}

	// TEMP: we need to remove the old nanobox-vpn just incase it is left over
	// we will not catch errors here because if it doesnt exist or it breaks it
	// should not stop us from creating the new nanobox-server
	service.Stop("nanobox-vpn")
	service.Remove("nanobox-vpn")

	// create the service this call is idempotent so we shouldnt need to check
	if err := service.Create("nanobox-server", []string{config.NanoboxPath(), "server"}); err != nil {
		return err
	}

	// start the service
	return Start()

}

// reExecPrivilageStart re-execs the current process with a privileged user
func reExecPrivilageStart() error {
	display.PauseTask()
	defer display.ResumeTask()

	display.PrintRequiresPrivilege("to start the server")

	cmd := fmt.Sprintf("%s env server start", config.NanoboxPath())

	if err := util.PrivilegeExec(cmd); err != nil {
		lumber.Error("server:reExecPrivilageStart:util.PrivilegeExec(%s): %s", cmd, err)
		return err
	}

	return nil
}
