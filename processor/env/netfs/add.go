package netfs

import (
	"fmt"
	"runtime"

	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/netfs"
)

// processEnvNetFSAdd ...
type processEnvNetFSAdd struct {
	control processor.ProcessControl
	path    string
}

func init() {
	processor.Register("env_netfs_add", envNetFSAddFn)
}

//
func envNetFSAddFn(control processor.ProcessControl) (processor.Processor, error) {
	envNetFSAdd := &processEnvNetFSAdd{control: control}
	return envNetFSAdd, envNetFSAdd.validateMeta()
}

func (envNetFSAdd processEnvNetFSAdd) Results() processor.ProcessControl {
	return envNetFSAdd.control
}

//
func (envNetFSAdd processEnvNetFSAdd) Process() error {

	// short-circuit if the entry already exist; we do this after we validate meta
	// because the meta is needed to determin the entry we're looking for
	if envNetFSAdd.entryExists() {
		return nil
	}

	// if we're not running as the privileged user, we need to re-exec with privilege
	if !util.IsPrivileged() {

		if err := envNetFSAdd.reExecPrivilege(); err != nil {
			return err
		}

		return nil
	}

	// add the netfs entry
	if err := envNetFSAdd.addEntry(); err != nil {
		return err
	}

	return nil
}

// validateMeta validates that the required metadata exists
func (envNetFSAdd *processEnvNetFSAdd) validateMeta() error {

	// set path (required) and ensure it's provided
	envNetFSAdd.path = envNetFSAdd.control.Meta["path"]
	if envNetFSAdd.path == "" {
		return fmt.Errorf("Missing required meta value 'path'")
	}

	return nil
}

// entryExists returns true if the entry already exists
func (envNetFSAdd *processEnvNetFSAdd) entryExists() bool {

	// if the entry exists just return
	if netfs.Exists(envNetFSAdd.path) {
		return true
	}

	return false
}

// addEntry adds the netfs entry into the /etc/exports
func (envNetFSAdd *processEnvNetFSAdd) addEntry() error {

	if err := netfs.Add(envNetFSAdd.path); err != nil {
		return err
	}

	return nil
}

// reExecPrivilege re-execs the current process with a privileged user
func (envNetFSAdd *processEnvNetFSAdd) reExecPrivilege() error {

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

	// call 'dev netfs add' with the original path and args; config.NanoboxPath() will be the
	// currently executing program, so this command will ultimately lead right back
	// here
	cmd := fmt.Sprintf("%s env netfs add %s", config.NanoboxPath(), envNetFSAdd.path)

	// if the escalated subprocess fails, we need to return error to stop the process
	if err := util.PrivilegeExec(cmd); err != nil {
		return err
	}

	return nil
}
