package dns

import (
	"fmt"

	"github.com/jcelliott/lumber"

	"github.com/nanobox-io/nanobox/commands/server"
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/display"
	"github.com/nanobox-io/nanobox/util/dns"
)

// RemoveAll removes all dns entries for an app
func RemoveAll(a *models.App) error {

	// shortcut if we dont have any entries for this app
	if len(dns.List(a.ID)) == 0 {

		return nil
	}

	// ensure we're running as the administrator for this
	if !util.IsPrivileged() {
		return reExecPrivilegedRemoveAll(a)
	}

	if err := dns.Remove(a.ID); err != nil {
		lumber.Error("dns:RemoveAll:dns.Remove(%s): %s", a.ID, err.Error())
		return util.ErrorAppend(err, "failed to remove all dns entries")
	}

	display.Info("\n%s removed all\n", display.TaskComplete)
	return nil
}

// reExecPrivilegedRemoveAll re-execs the current process with a privileged user
func reExecPrivilegedRemoveAll(a *models.App) error {
	display.PauseTask()
	defer display.ResumeTask()

	// display.PrintRequiresPrivilege("to modify host dns entries")

	// call 'dev dns add' with the original path and args
	resp, err := server.RunCommand("dns rm-all", []string{a.DisplayName()})

	// if the sudo'ed subprocess fails, we need to return error to stop the process
	if err != nil || resp == nil {
		lumber.Error("dns:reExecPrivilegedAdd:util.PrivilegeExec(dns add): %s", err)
		return err
	}

	if resp.ExitCode != 0 {
		lumber.Error("dns:reExecPrivilegedAdd:util.PrivilegeExec(dns add): %+v, %s", resp, err)
		return fmt.Errorf("bad exit code from server command(%d)", resp.ExitCode)		
	}

	fmt.Println(resp.Output)
	return nil
}
