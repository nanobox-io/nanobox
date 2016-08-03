package netfs

import (
	"fmt"
	"runtime"

	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/netfs"
)

// processEnvNetFSRemove ...
type processEnvNetFSRemove struct {
	control processor.ProcessControl
	path    string
}

func init() {
	processor.Register("env_netfs_remove", envNetFSRemoveFn)
}

//
func envNetFSRemoveFn(control processor.ProcessControl) (processor.Processor, error) {
	envNetFSRemove := &processEnvNetFSRemove{control: control}
	return envNetFSRemove, envNetFSRemove.validateMeta()
}

func (envNetFSRemove processEnvNetFSRemove) Results() processor.ProcessControl {
	return envNetFSRemove.control
}

//
func (envNetFSRemove processEnvNetFSRemove) Process() error {

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

// validateMeta validates that the required metadata exists
func (envNetFSRemove *processEnvNetFSRemove) validateMeta() error {

	// set path (required) and ensure it's provided
	envNetFSRemove.path = envNetFSRemove.control.Meta["path"]
	if envNetFSRemove.path == "" {
		return fmt.Errorf("Missing required meta value 'path'")
	}

	return nil
}

// entryExists returns true if the entry already exists
func (envNetFSRemove *processEnvNetFSRemove) entryExists() bool {

	// if the entry exists just return
	if netfs.Exists(envNetFSRemove.path) {
		return true
	}

	return false
}

// rmEntry rms the netfs entry into the /etc/exports
func (envNetFSRemove *processEnvNetFSRemove) rmEntry() error {

	// rm the entry into the /etc/exports file
	if err := netfs.Remove(envNetFSRemove.path); err != nil {
		return err
	}

	return nil
}

// reExecPrivilege re-execs the current process with a privileged user
func (envNetFSRemove *processEnvNetFSRemove) reExecPrivilege() error {

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
	cmd := fmt.Sprintf("%s env netfs rm %s", config.NanoboxPath(), envNetFSRemove.path)

	// if the sudo'ed subprocess fails, we need to return error to stop the process
	if err := util.PrivilegeExec(cmd); err != nil {
		return err
	}

	return nil
}
