package netfs

import (
	"fmt"
	"os"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/netfs"
)

// processDevNetFSAdd ...
type processDevNetFSAdd struct {
	control processor.ProcessControl
	app     models.App
	host    string
	path    string
}

//
func init() {
	processor.Register("dev_netfs_add", devNetFSAddFunc)
}

//
func devNetFSAddFunc(control processor.ProcessControl) (processor.Processor, error) {
	return processDevNetFSAdd{control: control}, nil
}

//
func (devNetFSAdd processDevNetFSAdd) Results() processor.ProcessControl {
	return devNetFSAdd.control
}

//
func (devNetFSAdd processDevNetFSAdd) Process() error {

	// validate we have all meta information needed and set "host" and "path"
	if err := devNetFSAdd.validateMeta(); err != nil {
		return err
	}

	// short-circuit if the entry already exist; we do this after we validate meta
	// because the meta is needed to determin the entry we're looking for
	if devNetFSAdd.entryExists() {
		return nil
	}

	// if we're not running as the privileged user, we need to re-exec with privilege
	if !util.IsPrivileged() {

		if err := devNetFSAdd.reExecPrivilege(); err != nil {
			return err
		}

		return nil
	}

	// add the netfs entry
	if err := devNetFSAdd.addEntry(); err != nil {
		return err
	}

	return nil
}

// validateMeta validates that the required metadata exists
func (devNetFSAdd processDevNetFSAdd) validateMeta() error {

	// set the host and path
	devNetFSAdd.host = devNetFSAdd.control.Meta["host"]
	devNetFSAdd.path = devNetFSAdd.control.Meta["path"]

	// ensure host and path are provided
	switch {
	case devNetFSAdd.host == "":
		return fmt.Errorf("Host is required")
	case devNetFSAdd.path == "":
		return fmt.Errorf("Path is required")
	}

	return nil
}

// entryExists returns true if the entry already exists
func (devNetFSAdd processDevNetFSAdd) entryExists() bool {

	// generate the entry
	entry := netfs.Entry(devNetFSAdd.host, devNetFSAdd.path)

	// if the entry exists just return
	if netfs.Exists(entry) {
		return true
	}

	return false
}

// addEntry adds the netfs entry into the /etc/exports
func (devNetFSAdd processDevNetFSAdd) addEntry() error {

	// generate the entry
	entry := netfs.Entry(devNetFSAdd.host, devNetFSAdd.path)

	// add the entry into the /etc/exports file
	if err := netfs.Add(entry); err != nil {
		return err
	}

	return nil
}

// reExecPrivilege re-execs the current process with a privileged user
func (devNetFSAdd processDevNetFSAdd) reExecPrivilege() error {

	// call 'dev netfs add' with the original path and args; os.Args[0] will be the
	// currently executing program, so this command will ultimately lead right back
	// here
	cmd := fmt.Sprintf("%s dev netfs add %s %s", os.Args[0], devNetFSAdd.host, devNetFSAdd.path)

	// if the sudo'ed subprocess fails, we need to return error to stop the process
	fmt.Println("Admin privileges are required to add entries to your exports file, your password may be requested...")
	if err := util.PrivilegeExec(cmd); err != nil {
		return err
	}

	return nil
}
