package netfs

import (
	"fmt"
	"os"

	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/netfs"
)

// processShareNetFSAdd ...
type processShareNetFSAdd struct {
	control processor.ProcessControl
	path    string
}

func init() {
	processor.Register("share_netfs_add", shareNetFSAddFn)
}

//
func shareNetFSAddFn(control processor.ProcessControl) (processor.Processor, error) {
	shareNetFSAdd := &processShareNetFSAdd{control: control}
	return shareNetFSAdd, shareNetFSAdd.validateMeta()
}

func (shareNetFSAdd processShareNetFSAdd) Results() processor.ProcessControl {
	return shareNetFSAdd.control
}

//
func (shareNetFSAdd processShareNetFSAdd) Process() error {

	// short-circuit if the entry already exist; we do this after we validate meta
	// because the meta is needed to determin the entry we're looking for
	if shareNetFSAdd.entryExists() {
		return nil
	}

	// if we're not running as the privileged user, we need to re-exec with privilege
	if !util.IsPrivileged() {

		if err := shareNetFSAdd.reExecPrivilege(); err != nil {
			return err
		}

		return nil
	}

	// add the netfs entry
	if err := shareNetFSAdd.addEntry(); err != nil {
		return err
	}

	return nil
}

// validateMeta validates that the required metadata exists
func (shareNetFSAdd *processShareNetFSAdd) validateMeta() error {

	// set path (required) and ensure it's provided
	shareNetFSAdd.path = shareNetFSAdd.control.Meta["path"]
	if shareNetFSAdd.path == "" {
		return fmt.Errorf("Missing required meta value 'path'")
	}

	return nil
}

// entryExists returns true if the entry already exists
func (shareNetFSAdd *processShareNetFSAdd) entryExists() bool {

	// if the entry exists just return
	if netfs.Exists(shareNetFSAdd.path) {
		return true
	}

	return false
}

// addEntry adds the netfs entry into the /etc/exports
func (shareNetFSAdd *processShareNetFSAdd) addEntry() error {

	if err := netfs.Add(shareNetFSAdd.path); err != nil {
		return err
	}

	return nil
}

// reExecPrivilege re-execs the current process with a privileged user
func (shareNetFSAdd *processShareNetFSAdd) reExecPrivilege() error {

	// call 'dev netfs add' with the original path and args; os.Args[0] will be the
	// currently executing program, so this command will ultimately lead right back
	// here
	cmd := fmt.Sprintf("%s dev netfs add %s", os.Args[0], shareNetFSAdd.path)

	// if the sudo'ed subprocess fails, we need to return error to stop the process
	fmt.Println("Admin privileges are required to add entries to your exports file, your password may be requested...")
	if err := util.PrivilegeExec(cmd); err != nil {
		return err
	}

	return nil
}
