package dns

import (
	"fmt"

	"github.com/jcelliott/lumber"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/config"
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
		return fmt.Errorf("failed to remove all dns entries: %s", err.Error())
	}

	return nil
}

// reExecPrivilegedRemoveAll re-execs the current process with a privileged user
func reExecPrivilegedRemoveAll(a *models.App) error {
	display.PauseTask()
	defer display.ResumeTask()

	display.PrintRequiresPrivilege("to modify host dns entries")

	// call 'dev dns add' with the original path and args
	cmd := fmt.Sprintf("%s %s dns rm-all", config.NanoboxPath(), a.Name)

	// if the sudo'ed subprocess fails, we need to return error to stop the process
	if err := util.PrivilegeExec(cmd); err != nil {
		lumber.Error("dns:reExecPrivilegedRemoveAll:util.PrivilegeExec(%s): %s", cmd, err)
		return err
	}

	return nil
}
