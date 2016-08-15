package netfs

import (
	"fmt"
	"runtime"

	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/netfs"
)

// Remove ...
type Remove struct {
	Path string
}

//
func (envNetFSRemove Remove) Run() error {

	// short-circuit if the entry already exist; we do this after we validate meta
	// because the meta is needed to determin the entry we're looking for
	if envNetFSRemove.entryExists() {
		return nil
	}

	// if we're not running as the privileged user, we need to re-exec with privilege
	if !util.IsPrivileged() {

		if err := envNetFSRemove.reExecPrivilege(); err != nil {
			return err
		}

		return nil
	}

	// rm the netfs entry
	if err := envNetFSRemove.rmEntry(); err != nil {
		return err
	}

	return nil
}

// entryExists returns true if the entry already exists
func (envNetFSRemove *Remove) entryExists() bool {

	// if the entry exists just return
	if netfs.Exists(envNetFSRemove.Path) {
		return true
	}

	return false
}

// rmEntry rms the netfs entry into the /etc/exports
func (envNetFSRemove *Remove) rmEntry() error {

	// rm the entry into the /etc/exports file
	if err := netfs.Remove(envNetFSRemove.Path); err != nil {
		return err
	}

	return nil
}

// reExecPrivilege re-execs the current process with a privileged user
func (envNetFSRemove *Remove) reExecPrivilege() error {

	if runtime.GOOS == "windows" {
		fmt.Println("Administrator privileges are required to modify network shares.")
		fmt.Println("Another window will be opened as the Administrator and permission may be requested.")

		// block here until the user hits enter. It's not ideal, but we need to make
		// sure they see the new window open if it requests their password.
		fmt.Println("Enter to continue:")
		var input string
		fmt.Scanln(&input)
	} else {
		fmt.Println("Root privileges are required to modify your exports file, your password may be requested...")
	}

	// call 'env netfs rm' with the original path and args; config.NanoboxPath() will be the
	// currently executing program, so this command will ultimately lead right back
	// here
	cmd := fmt.Sprintf("%s env netfs rm %s", config.NanoboxPath(), envNetFSRemove.Path)

	// if the sudo'ed subprocess fails, we need to return error to stop the process
	if err := util.PrivilegeExec(cmd); err != nil {
		return err
	}

	return nil
}
