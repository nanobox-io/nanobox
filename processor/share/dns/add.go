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

// processShareDNSAdd ...
type processShareDNSAdd struct {
	control processor.ProcessControl
	app     models.App
	name    string
}

func init() {
	processor.Register("share_dns_add", shareDNSAddFn)
}

// shareDNSAddFn creates a processShareDNSAdd and validates the meta data in the control
func shareDNSAddFn(control processor.ProcessControl) (processor.Processor, error) {
	shareDNSAdd := &processShareDNSAdd{control: control}
	return shareDNSAdd, shareDNSAdd.validateMeta()
}

func (shareDNSAdd processShareDNSAdd) Results() processor.ProcessControl {
	return shareDNSAdd.control
}

//
func (shareDNSAdd processShareDNSAdd) Process() error {

	// load the current "app"
	if err := shareDNSAdd.loadApp(); err != nil {
		return err
	}

	// short-circuit if the entries already exist; we do this after we validate meta
	// and load the app because both of those are needed to determin the entry we're
	// looking for
	if shareDNSAdd.entriesExist() {
		return nil
	}

	// if we're not running as the privileged user, we need to re-exec with privilege
	if !util.IsPrivileged() {
		return shareDNSAdd.reExecPrivilege()
	}

	// add the DNS entries
	return shareDNSAdd.addEntries()
}

// validateMeta validates that the required metadata exists
func (shareDNSAdd *processShareDNSAdd) validateMeta() error {

	// set name (required) and ensure it's provided
	shareDNSAdd.name = shareDNSAdd.control.Meta["name"]
	if shareDNSAdd.name == "" {
		return fmt.Errorf("Missing required meta value 'name'")
	}

	return nil
}

// loadApp loads the app from the db
func (shareDNSAdd *processShareDNSAdd) loadApp() error {

	//
	key := fmt.Sprintf("%s_%s", config.AppName(), shareDNSAdd.control.Env)
	return data.Get("apps", key, &shareDNSAdd.app)
}

// entriesExist returns true if both entries already exist
func (shareDNSAdd *processShareDNSAdd) entriesExist() bool {

	// fetch the IP 
	// env in dev is used in the dev container
	// env in sim is used for portal
	envIP := shareDNSAdd.app.GlobalIPs["env"]

	// generate the entries
	env := dns.Entry(envIP, shareDNSAdd.name, shareDNSAdd.control.Env)

	// if the entries dont exist just return
	return dns.Exists(env)
}

// addEntries adds the dev and sim entries into the host dns
func (shareDNSAdd *processShareDNSAdd) addEntries() error {

	// fetch the IPs
	envIP := shareDNSAdd.app.GlobalIPs["env"]

	// generate the entries
	env := dns.Entry(envIP, shareDNSAdd.name, shareDNSAdd.control.Env)

	// add the 'sim' DNS entry into the /etc/hosts file
	return dns.Add(env)
}

// reExecPrivilege re-execs the current process with a privileged user
func (shareDNSAdd *processShareDNSAdd) reExecPrivilege() error {

	// call 'dev dns add' with the original path and args; os.Args[0] will be the
	// currently executing program, so this command will ultimately lead right back
	// here
	cmd := fmt.Sprintf("%s %s dns add %s", os.Args[0], shareDNSAdd.control.Env, shareDNSAdd.name)

	// if the sudo'ed subprocess fails, we need to return error to stop the process
	fmt.Println("Admin privileges are required to add DNS entries to your hosts file, your password may be requested...")
	if err := util.PrivilegeExec(cmd); err != nil {
		return err
	}

	return nil
}
