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

var AppSetup func(envModel *models.Env, appModel *models.App, name string) error

// Add adds a dns entry to the local hosts file
func Add(envModel *models.Env, appModel *models.App, name string) error {

	if err := AppSetup(envModel, appModel, appModel.Name); err != nil {
		return util.ErrorAppend(err, "failed to setup app")
	}

	// fetch the IP
	// env in dev is used in the dev container
	// env in sim is used for portal
	envIP := appModel.LocalIPs["env"]

	// generate the dns entry
	entry := dns.Entry(envIP, name, appModel.ID)

	// short-circuit if this entry already exists
	if dns.Exists(entry) {
		return nil
	}

	// ensure we're running as the administrator for this
	if !util.IsPrivileged() {
		return reExecPrivilegedAdd(appModel, name)
	}

	// add the entry
	if err := dns.Add(entry); err != nil {
		lumber.Error("dns:Add:dns.Add(%s): %s", entry, err.Error())
		return util.ErrorAppend(err, "unable to add dns entry")
	}

	display.Info("\n%s %s added\n", display.TaskComplete, name)

	return nil
}

// reExecPrivilegedAdd re-execs the current process with a privileged user
func reExecPrivilegedAdd(appModel *models.App, name string) error {
	display.PauseTask()
	defer display.ResumeTask()

	// display.PrintRequiresPrivilege("to modify host dns entries")

	// // call 'dev dns add' with the original path and args
	// cmd := fmt.Sprintf("%s dns add %s %s", config.NanoboxPath(), appModel.DisplayName(), name)

	resp, err := server.RunCommand("dns add", []string{appModel.DisplayName(), name})

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
