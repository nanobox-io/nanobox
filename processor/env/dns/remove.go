package dns

import (
	"fmt"
	"os"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/data"
	"github.com/nanobox-io/nanobox/util/dns"
)

// processEnvDNSRemove ...
type processEnvDNSRemove struct {
	control processor.ProcessControl
	app     models.App
	name    string
}

func init() {
	processor.Register("env_dns_remove", envDNSRemoveFn)
}

// devDNSRemveFn creates a processEnvDNSRemove and validates the meta in the control
func envDNSRemoveFn(control processor.ProcessControl) (processor.Processor, error) {
	envDNSRemove := &processEnvDNSRemove{control: control}
	return envDNSRemove, envDNSRemove.validateMeta()
}

func (envDNSRemove processEnvDNSRemove) Results() processor.ProcessControl {
	return envDNSRemove.control
}

//
func (envDNSRemove processEnvDNSRemove) Process() error {

	// load the current "app"
	if err := envDNSRemove.loadApp(); err != nil {
		return err
	}

	// short-circuit if the entries dont exist; we do this after we validate meta
	// and load the app because both of those are needed to determin the entry we're
	// looking for
	if !envDNSRemove.entriesExist() {
		return nil
	}

	// if we're not running as the privileged user, we need to re-exec with privilege
	if !util.IsPrivileged() {
		return envDNSRemove.reExecPrivilege()
	}

	// remove the DNS entries
	return envDNSRemove.removeEntries()
}

// validateMeta validates that the required metadata exists
func (envDNSRemove *processEnvDNSRemove) validateMeta() error {

	// set name (required) and ensure it's provided
	envDNSRemove.name = envDNSRemove.control.Meta["name"]
	if envDNSRemove.name == "" {
		return fmt.Errorf("Missing required meta value 'name'")
	}

	return nil
}

// loadApp loads the app from the db
func (envDNSRemove *processEnvDNSRemove) loadApp() error {

	//
	key := fmt.Sprintf("%s_%s", config.AppID(), envDNSRemove.control.Env)
	return data.Get("apps", key, &envDNSRemove.app)
}

// entriesExist returns true if both entries already exist
func (envDNSRemove *processEnvDNSRemove) entriesExist() bool {

	// fetch the IPs
	envIP := envDNSRemove.app.GlobalIPs["env"]

	// generate the entries
	env := dns.Entry(envIP, envDNSRemove.name, fmt.Sprintf("%s_%s", config.AppID(), envDNSRemove.control.Env))

	return dns.Exists(env)
}

// removeEntries removes the "dev" and "env" entries into the host dns
func (envDNSRemove *processEnvDNSRemove) removeEntries() error {

	// fetch the IPs
	envIP := envDNSRemove.app.GlobalIPs["env"]

	// generate the entries
	env := dns.Entry(envIP, envDNSRemove.name, fmt.Sprintf("%s_%s", config.AppID(), envDNSRemove.control.Env))

	// remove the DNS entry from the /etc/hosts file
	return dns.Remove(env)
}

// reExecPrivilege re-execs the current process with a privileged user
func (envDNSRemove *processEnvDNSRemove) reExecPrivilege() error {

	// call 'dev dns rm' with the original path and args; os.Args[0] will be the
	// currently executing program, so this command will ultimately lead right back
	// here
	cmd := fmt.Sprintf("%s %s dns rm %s", os.Args[0], envDNSRemove.control.Env, envDNSRemove.name)

	// if the sudo'ed subprocess fails, we need to return error to stop the process
	fmt.Println("Admin privileges are required to remove DNS entries from your hosts file...")
	return util.PrivilegeExec(cmd)
}
