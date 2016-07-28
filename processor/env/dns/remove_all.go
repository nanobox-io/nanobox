package dns

import (
	"fmt"
	"os"
	"runtime"

	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/dns"
)

// processEnvDNSRemoveAll ...
type processEnvDNSRemoveAll struct {
	control processor.ProcessControl
}

func init() {
	processor.Register("env_dns_remove_all", envDNSRemoveAllFn)
}

// devDNSRemoveAllFn creates a processEnvDNSRemoveAll and validates the meta in the control
func envDNSRemoveAllFn(control processor.ProcessControl) (processor.Processor, error) {
	envDNSRemoveAll := &processEnvDNSRemoveAll{control: control}
	return envDNSRemoveAll, envDNSRemoveAll.validateMeta()
}

func (envDNSRemoveAll processEnvDNSRemoveAll) Results() processor.ProcessControl {
	return envDNSRemoveAll.control
}

//
func (envDNSRemoveAll processEnvDNSRemoveAll) Process() error {

	// if we're not running as the privileged user, we need to re-exec with privilege
	if !util.IsPrivileged() {
		return envDNSRemoveAll.reExecPrivilege()
	}

	// remove all the dns's that contain my app name
	return dns.Remove(envDNSRemoveAll.control.Meta["app_name"])
}

// validateMeta validates that the required metadata exists
func (envDNSRemoveAll *processEnvDNSRemoveAll) validateMeta() error {

	if envDNSRemoveAll.control.Meta["app_name"] == "" {
		envDNSRemoveAll.control.Meta["app_name"] = fmt.Sprintf("%s_%s", config.AppID(), envDNSRemoveAll.control.Env)
	}

	return nil
}

// reExecPrivilege re-execs the current process with a privileged user
func (envDNSRemoveAll *processEnvDNSRemoveAll) reExecPrivilege() error {

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
	cmd := fmt.Sprintf("%s %s dns rm-all", os.Args[0], envDNSRemoveAll.control.Env)

	// if the sudo'ed subprocess fails, we need to return error to stop the process
	return util.PrivilegeExec(cmd)
}
