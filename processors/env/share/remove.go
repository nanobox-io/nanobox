package share

import (
	"fmt"

	"github.com/jcelliott/lumber"

	"github.com/nanobox-io/nanobox/commands/server"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/display"
	"github.com/nanobox-io/nanobox/util/provider/share"
)

// Remove removes a share share from the workstation
func Remove(path string) error {

	// short-circuit if the entry doesn't exist
	if !share.Exists(path) {
		return nil
	}

	// // if we're not running as the privileged user, we need to re-exec with privilege
	// if !util.IsPrivileged() {
	// 	return reExecPrivilegedRemove(path)
	// }

	// rm the share entry
	if err := share.Remove(path); err != nil {
		lumber.Error("share:Add:share.Remove(%s): %s", path, err.Error())
		return util.ErrorAppend(err, "failed to remove share share")
	}

	return nil
}

// reExecPrivilegedRemove re-execs the current process with a privileged user
func reExecPrivilegedRemove(path string) error {
	display.PauseTask()
	defer display.ResumeTask()

	// display.PrintRequiresPrivilege("to modify network shares")

	// call 'dev share add' with the original path and args; config.NanoboxPath() will be the
	// currently executing program, so this command will ultimately lead right back
	// here
	resp, err := server.RunCommand("share rm", []string{path})

	// if the sudo'ed subprocess fails, we need to return error to stop the process
	if err != nil || resp == nil {
		lumber.Error("share:reExecPrivilegedAdd:util.PrivilegeExec(share rm): %s", err)
		return err
	}

	if resp.ExitCode != 0 {
		lumber.Error("share:reExecPrivilegedAdd:util.PrivilegeExec(share rm): %+v, %s", resp, err)
		return fmt.Errorf("bad exit code from server command(%d)", resp.ExitCode)
	}

	fmt.Printf(resp.Output)

	return nil
}
