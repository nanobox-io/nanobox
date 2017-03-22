package dns

import (
	"fmt"

	"github.com/jcelliott/lumber"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/display"
	"github.com/nanobox-io/nanobox/commands/server"
	"github.com/nanobox-io/nanobox/util/dns"
)

// Remove removes a dns entry from the local hosts file
func Remove(a *models.App, name string) error {
	// fetch the IP
	// env in dev is used in the dev container
	// env in sim is used for portal
	envIP := a.LocalIPs["env"]

	// generate the dns entry
	entry := dns.Entry(envIP, name, a.ID)

	// short-circuit if this entry doesn't exist
	if !dns.Exists(entry) {
		return nil
	}

	// ensure we're running as the administrator for this
	if !util.IsPrivileged() {
		return reExecPrivilegedRemove(a, name)
	}

	// remove the entry
	if err := dns.Remove(entry); err != nil {
		lumber.Error("dns:Remove:dns.Remove(%s): %s", entry, err.Error())
		return util.ErrorAppend(err, "unable to add dns entry: %s")
	}

	display.Info("\n%s %s removed\n", display.TaskComplete, name)

	return nil
}

// reExecPrivilegedRemove re-execs the current process with a privileged user
func reExecPrivilegedRemove(a *models.App, name string) error {
	display.PauseTask()
	defer display.ResumeTask()

	// display.PrintRequiresPrivilege("to modify host dns entries")

	// call 'dev dns add' with the original path and args
	resp, err := server.RunCommand("dns rm", []string{a.DisplayName(), name})

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
