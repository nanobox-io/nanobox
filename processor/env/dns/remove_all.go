package dns

import (
	"fmt"
	"os"

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
		envDNSRemoveAll.control.Meta["app_name"] = fmt.Sprintf("%s_%s", config.AppName(), envDNSRemoveAll.control.Env)
	}

	return nil
}

// reExecPrivilege re-execs the current process with a privileged user
func (envDNSRemoveAll *processEnvDNSRemoveAll) reExecPrivilege() error {

	// call rm-all
	cmd := fmt.Sprintf("%s %s dns rm-all", os.Args[0], envDNSRemoveAll.control.Env)

	// if the sudo'ed subprocess fails, we need to return error to stop the process
	fmt.Println("Admin privileges are required to remove DNS entries from your hosts file...")
	return util.PrivilegeExec(cmd)
}
