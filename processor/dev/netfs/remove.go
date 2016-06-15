package netfs

import (
	"fmt"
	"os"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/netfs"
)

// processDevNetFSRemove ...
type processDevNetFSRemove struct {
	control processor.ProcessControl
	app     models.App
	host    string
	path    string
}

//
func init() {
	processor.Register("dev_netfs_remove", devNetFSRemoveFn)
}

//
func devNetFSRemoveFn(control processor.ProcessControl) (processor.Processor, error) {
	return processDevNetFSRemove{control: control}, nil
}

//
func (devNetFSRemove processDevNetFSRemove) Results() processor.ProcessControl {
	return devNetFSRemove.control
}

//
func (devNetFSRemove processDevNetFSRemove) Process() error {

	// validate we have all meta information needed and set "host" and "path"
	if err := devNetFSRemove.validateMeta(); err != nil {
		return err
	}

	// short-circuit if the entry doesnt exist; we do this after we validate meta
	// because the meta is needed to determin the entry we're looking for
	if !devNetFSRemove.entryExists() {
		return nil
	}

	// if we're not running as the privileged user, we need to re-exec with privilege
	if !util.IsPrivileged() {

		if err := devNetFSRemove.reExecPrivilege(); err != nil {
			return err
		}

		return nil
	}

	// remove the netfs entry
	if err := devNetFSRemove.removeEntry(); err != nil {
		return err
	}

	return nil
}

// validateMeta validates that the required metadata exists
func (devNetFSRemove processDevNetFSRemove) validateMeta() error {

	// set the host and path
	devNetFSRemove.host = devNetFSRemove.control.Meta["host"]
	devNetFSRemove.path = devNetFSRemove.control.Meta["path"]

	// ensure host and path are provided
	switch {
	case devNetFSRemove.host == "":
		return fmt.Errorf("Host is required")
	case devNetFSRemove.path == "":
		return fmt.Errorf("Path is required")
	}

	return nil
}

// entryExists returns true if the entry already exists
func (devNetFSRemove processDevNetFSRemove) entryExists() bool {

	// generate the entry
	entry := netfs.Entry(devNetFSRemove.host, devNetFSRemove.path)

	// if the entry exists just return
	if netfs.Exists(entry) {
		return true
	}

	return false
}

// removeEntry removes the netfs entry from the /etc/exports
func (devNetFSRemove processDevNetFSRemove) removeEntry() error {

	// generate the entry
	entry := netfs.Entry(devNetFSRemove.host, devNetFSRemove.path)

	// remove the entry from the /etc/exports file
	if err := netfs.Remove(entry); err != nil {
		return err
	}

	return nil
}

// reExecPrivilege re-execs the current process with a privileged user
func (devNetFSRemove processDevNetFSRemove) reExecPrivilege() error {

	// call 'dev netfs rm' with the original path and args; os.Args[0] will be the
	// currently executing program, so this command will ultimately lead right back
	// here
	cmd := fmt.Sprintf("%s dev netfs rm %s %s", os.Args[0], devNetFSRemove.host, devNetFSRemove.path)

	// if the sudo'ed subprocess fails, we need to return error to stop the process
	fmt.Println("Admin privileges are required to remove entries from your exports file, your password may be requested...")
	if err := util.PrivilegeExec(cmd); err != nil {
		return err
	}

	return nil
}
