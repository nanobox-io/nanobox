package share

import (
	"fmt"

	"github.com/jcelliott/lumber"

	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/display"
	"github.com/nanobox-io/nanobox/util/provider/share"
)

// Add adds a share share to the workstation
func Add(path string) error {

	// short-circuit if the entry already exist
	if share.Exists(path) {
		return nil
	}

	// if we're not running as the privileged user, we need to re-exec with privilege
	if !util.IsPrivileged() {

		return reExecPrivilegedAdd(path)
	}

	// add the share entry
	if err := share.Add(path); err != nil {
		lumber.Error("share:Add:share.Add(%s): %s", path, err.Error())
		return util.ErrorAppend(err, "failed to add share")
	}

	return nil
}

// reExecPrivilegedAdd re-execs the current process with a privileged user
func reExecPrivilegedAdd(path string) error {
	display.PauseTask()
	defer display.ResumeTask()

	display.PrintRequiresPrivilege("to modify network shares")

	// call 'dev share add' with the original path and args; config.NanoboxPath() will be the
	// currently executing program, so this command will ultimately lead right back
	// here
	cmd := fmt.Sprintf("%s env share add \"%s\"", config.NanoboxPath(), path)

	// if the escalated subprocess fails, we need to return error to stop the process
	if err := util.PrivilegeExec(cmd); err != nil {
		lumber.Error("share:reExecPrivilegedAdd:util.PrivilegeExec(%s): %s", cmd, err)
		return err
	}

	return nil
}
