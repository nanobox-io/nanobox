package processor

import (
	"fmt"
	"runtime"

	"github.com/jcelliott/lumber"

	"github.com/nanobox-io/nanobox/processor/provider"
	"github.com/nanobox-io/nanobox/util/display"
	"github.com/nanobox-io/nanobox/util/dns"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/config"

)

// Destroy ...
type Destroy struct {
}

//
func (destroy Destroy) Run() error {
	display.OpenContext("Destroying nanobox system")
	defer display.CloseContext()

	// only run the provider destroy if not privilaged
	if !util.IsPrivileged() {
		providerDestroy := provider.Destroy{}
		// run a provider destroy
		if err := providerDestroy.Run(); err != nil {
			return err
		}
	}

	// stop here if no dns's were put in by nanobox
	if len(dns.List("by nanobox")) == 0 {
		return nil
	}

	// if we're not running as the privileged user, we need to re-exec with privilege
	if !util.IsPrivileged() {
		return destroy.reExecPrivilege()
	}

	return dns.Remove("by nanobox")
}


// reExecPrivilege re-execs the current process with a privileged user
func (destroy *Destroy) reExecPrivilege() error {

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
	cmd := fmt.Sprintf("%s destroy", config.NanoboxPath())

	// if the sudo'ed subprocess fails, we need to return error to stop the process
	if err := util.PrivilegeExec(cmd); err != nil {
		lumber.Error("dns:RemoveAll:reExecPrivilege:util.PrivilegeExec(%s): %s", cmd, err)
		return err
	}
	return nil
}

