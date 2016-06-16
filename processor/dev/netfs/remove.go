package netfs

import (
	"fmt"
	"os"

	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/netfs"
)

// processDevNetFSRemove ...
type processDevNetFSRemove struct {
	control processor.ProcessControl
	path    string
}

//
func init() {
	processor.Register("dev_netfs_remove", devNetFSRemoveFn)
}

//
func devNetFSRemoveFn(control processor.ProcessControl) (processor.Processor, error) {
	devNetFSRemove := processDevNetFSRemove{control: control}
	return devNetFSRemove, devNetFSRemove.validateMeta()
}

//
func (devNetFSRemove processDevNetFSRemove) Results() processor.ProcessControl {
	return devNetFSRemove.control
}

//
func (devNetFSRemove processDevNetFSRemove) Process() error {

	// short-circuit if the entry already exist; we do this after we validate meta
	// because the meta is needed to determin the entry we're looking for
	if devNetFSRemove.entryExists() {
		return nil
	}

	// if we're not running as the privileged user, we need to re-exec with privilege
	if !util.IsPrivileged() {

		if err := devNetFSRemove.reExecPrivilege(); err != nil {
			return err
		}

		return nil
	}

	// rm the netfs entry
	if err := devNetFSRemove.rmEntry(); err != nil {
		return err
	}

	return nil
}

// validateMeta validates that the required metadata exists
func (devNetFSRemove *processDevNetFSRemove) validateMeta() error {

	// set the host and path
	devNetFSRemove.path = devNetFSRemove.control.Meta["path"]

	if devNetFSRemove.path == "" {
		return fmt.Errorf("Path is required")
	}

	return nil
}

// entryExists returns true if the entry already exists
func (devNetFSRemove *processDevNetFSRemove) entryExists() bool {

	// if the entry exists just return
	if netfs.Exists(devNetFSRemove.path) {
		return true
	}

	return false
}

// rmEntry rms the netfs entry into the /etc/exports
func (devNetFSRemove *processDevNetFSRemove) rmEntry() error {

	// rm the entry into the /etc/exports file
	if err := netfs.Remove(devNetFSRemove.path); err != nil {
		return err
	}

	return nil
}

// reExecPrivilege re-execs the current process with a privileged user
func (devNetFSRemove *processDevNetFSRemove) reExecPrivilege() error {

	// call 'dev netfs rm' with the original path and args; os.Args[0] will be the
	// currently executing program, so this command will ultimately lead right back
	// here
	cmd := fmt.Sprintf("%s dev netfs rm %s", os.Args[0], devNetFSRemove.path)

	// if the sudo'ed subprocess fails, we need to return error to stop the process
	fmt.Println("Admin privileges are required to remove entries from your exports file, your password may be requested...")
	if err := util.PrivilegeExec(cmd); err != nil {
		return err
	}

	return nil
}
