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

// processShareDNSRemove ...
type processShareDNSRemove struct {
	control processor.ProcessControl
	app     models.App
	name    string
}

func init() {
	processor.Register("share_dns_remove", shareDNSRemoveFn)
}

// devDNSRemveFn creates a processShareDNSRemove and validates the meta in the control
func shareDNSRemoveFn(control processor.ProcessControl) (processor.Processor, error) {
	shareDNSRemove := &processShareDNSRemove{control: control}
	return shareDNSRemove, shareDNSRemove.validateMeta()
}

func (shareDNSRemove processShareDNSRemove) Results() processor.ProcessControl {
	return shareDNSRemove.control
}

//
func (shareDNSRemove processShareDNSRemove) Process() error {

	// load the current "app"
	if err := shareDNSRemove.loadApp(); err != nil {
		return err
	}

	// short-circuit if the entries dont exist; we do this after we validate meta
	// and load the app because both of those are needed to determin the entry we're
	// looking for
	if !shareDNSRemove.entriesExist() {
		return nil
	}

	// if we're not running as the privileged user, we need to re-exec with privilege
	if !util.IsPrivileged() {
		return shareDNSRemove.reExecPrivilege()
	}

	// remove the DNS entries
	return shareDNSRemove.removeEntries()
}

// validateMeta validates that the required metadata exists
func (shareDNSRemove *processShareDNSRemove) validateMeta() error {

	// set name (required) and ensure it's provided
	shareDNSRemove.name = shareDNSRemove.control.Meta["name"]
	if shareDNSRemove.name == "" {
		return fmt.Errorf("Missing required meta value 'name'")
	}

	return nil
}

// loadApp loads the app from the db
func (shareDNSRemove *processShareDNSRemove) loadApp() error {

	//
	key := fmt.Sprintf("%s_%s", config.AppName(), shareDNSRemove.control.Env)
	return data.Get("apps", key, &shareDNSRemove.app)
}

// entriesExist returns true if both entries already exist
func (shareDNSRemove *processShareDNSRemove) entriesExist() bool {

	// fetch the IPs
	envIP := shareDNSRemove.app.GlobalIPs["env"]

	// generate the entries
	env := dns.Entry(envIP, shareDNSRemove.name, shareDNSRemove.control.Env)

	return dns.Exists(env)
}

// removeEntries removes the "dev" and "env" entries into the host dns
func (shareDNSRemove *processShareDNSRemove) removeEntries() error {

	// fetch the IPs
	envIP := shareDNSRemove.app.GlobalIPs["env"]

	// generate the entries
	env := dns.Entry(envIP, shareDNSRemove.name, shareDNSRemove.control.Env)

	// remove the DNS entry from the /etc/hosts file
	return dns.Remove(env)
}

// reExecPrivilege re-execs the current process with a privileged user
func (shareDNSRemove *processShareDNSRemove) reExecPrivilege() error {

	// call 'dev dns rm' with the original path and args; os.Args[0] will be the
	// currently executing program, so this command will ultimately lead right back
	// here
	cmd := fmt.Sprintf("%s %s dns rm %s", os.Args[0], shareDNSRemove.control.Env, shareDNSRemove.name)

	// if the sudo'ed subprocess fails, we need to return error to stop the process
	fmt.Println("Admin privileges are required to remove DNS entries from your hosts file, your password may be requested...")
	return util.PrivilegeExec(cmd)
}
