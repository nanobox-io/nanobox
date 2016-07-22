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

// processEnvDNSAdd ...
type processEnvDNSAdd struct {
	control processor.ProcessControl
	app     models.App
	name    string
}

func init() {
	processor.Register("env_dns_add", envDNSAddFn)
}

// envDNSAddFn creates a processEnvDNSAdd and validates the meta data in the control
func envDNSAddFn(control processor.ProcessControl) (processor.Processor, error) {
	envDNSAdd := &processEnvDNSAdd{control: control}
	return envDNSAdd, envDNSAdd.validateMeta()
}

func (envDNSAdd processEnvDNSAdd) Results() processor.ProcessControl {
	return envDNSAdd.control
}

//
func (envDNSAdd processEnvDNSAdd) Process() error {

	// load the current "app"
	if err := envDNSAdd.loadApp(); err != nil {
		return err
	}

	// short-circuit if the entries already exist; we do this after we validate meta
	// and load the app because both of those are needed to determin the entry we're
	// looking for
	if envDNSAdd.entriesExist() {
		return nil
	}

	// if we're not running as the privileged user, we need to re-exec with privilege
	if !util.IsPrivileged() {
		return envDNSAdd.reExecPrivilege()
	}

	// add the DNS entries
	return envDNSAdd.addEntries()
}

// validateMeta validates that the required metadata exists
func (envDNSAdd *processEnvDNSAdd) validateMeta() error {

	// set name (required) and ensure it's provided
	envDNSAdd.name = envDNSAdd.control.Meta["name"]
	if envDNSAdd.name == "" {
		return fmt.Errorf("Missing required meta value 'name'")
	}

	return nil
}

// loadApp loads the app from the db
func (envDNSAdd *processEnvDNSAdd) loadApp() error {

	//
	key := fmt.Sprintf("%s_%s", config.AppName(), envDNSAdd.control.Env)
	return data.Get("apps", key, &envDNSAdd.app)
}

// entriesExist returns true if both entries already exist
func (envDNSAdd *processEnvDNSAdd) entriesExist() bool {

	// fetch the IP
	// env in dev is used in the dev container
	// env in sim is used for portal
	envIP := envDNSAdd.app.GlobalIPs["env"]

	// generate the entries
	env := dns.Entry(envIP, envDNSAdd.name, fmt.Sprintf("%s_%s", config.AppName(), envDNSAdd.control.Env))

	// if the entries dont exist just return
	return dns.Exists(env)
}

// addEntries adds the dev and sim entries into the host dns
func (envDNSAdd *processEnvDNSAdd) addEntries() error {

	// fetch the IPs
	envIP := envDNSAdd.app.GlobalIPs["env"]

	// generate the entries
	env := dns.Entry(envIP, envDNSAdd.name, fmt.Sprintf("%s_%s", config.AppName(), envDNSAdd.control.Env))

	// add the 'sim' DNS entry into the /etc/hosts file
	return dns.Add(env)
}

// reExecPrivilege re-execs the current process with a privileged user
func (envDNSAdd *processEnvDNSAdd) reExecPrivilege() error {

	// call 'dev dns add' with the original path and args; os.Args[0] will be the
	// currently executing program, so this command will ultimately lead right back
	// here
	cmd := fmt.Sprintf("%s %s dns add %s", os.Args[0], envDNSAdd.control.Env, envDNSAdd.name)

	// if the sudo'ed subprocess fails, we need to return error to stop the process
	fmt.Println("Admin privileges are required to add DNS entries to your hosts file, your password may be requested...")
	if err := util.PrivilegeExec(cmd); err != nil {
		return err
	}

	return nil
}
