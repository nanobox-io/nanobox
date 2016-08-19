package dns

import (
	"fmt"
	"runtime"

	"github.com/jcelliott/lumber"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/dns"
)

// RemoveAll ...
type RemoveAll struct {
	App models.App
}

//
func (removeAll RemoveAll) Run() error {

	// if we're not running as the privileged user, we need to re-exec with privilege
	if !util.IsPrivileged() {
		return removeAll.reExecPrivilege()
	}

	// remove all the dns's that contain my app name
	return dns.Remove(removeAll.App.ID)
}

// reExecPrivilege re-execs the current process with a privileged user
func (removeAll *RemoveAll) reExecPrivilege() error {

	if runtime.GOOS == "windows" {
		fmt.Println("Administrator privileges are required to modify host dns entries.")
		fmt.Println("Another window will be opened as the Administrator and permission may be requested.")

		// block here until the user hits enter. It's not ideal, but we need to make
		// sure they see the new window open if it requests their password.
		fmt.Println("Enter to continue:")
		var input string
		fmt.Scanln(&input)
	} else {
		fmt.Println("Root privileges are required to modify your hosts file, your password may be requested...")
	}

	// call rm-all
	cmd := fmt.Sprintf("%s %s dns rm-all", config.NanoboxPath(), removeAll.App.Name)

	// if the sudo'ed subprocess fails, we need to return error to stop the process
	if err := util.PrivilegeExec(cmd); err != nil {
		lumber.Error("dns:RemoveAll:reExecPrivilege:util.PrivilegeExec(%s): %s", cmd, err)
		return err
	}
	return nil
}
