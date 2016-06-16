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

// processDevDNSAdd ...
type processDevDNSAdd struct {
	control processor.ProcessControl
	app     models.App
	name    string
}

//
func init() {
	processor.Register("dev_dns_add", devDNSAddFn)
}

//
func devDNSAddFn(control processor.ProcessControl) (processor.Processor, error) {
	return processDevDNSAdd{control: control}, nil
}

//
func (devDNSAdd processDevDNSAdd) Results() processor.ProcessControl {
	return devDNSAdd.control
}

//
func (devDNSAdd processDevDNSAdd) Process() error {

	// validate we have all meta information needed and set "name"
	if err := devDNSAdd.validateMeta(); err != nil {
		return err
	}

	// load the current "app"
	if err := devDNSAdd.loadApp(); err != nil {
		return err
	}

	// short-circuit if the entries already exist; we do this after we validate meta
	// and load the app because both of those are needed to determin the entry we're
	// looking for
	if devDNSAdd.entriesExist() {
		return nil
	}

	// if we're not running as the privileged user, we need to re-exec with privilege
	if !util.IsPrivileged() {

		if err := devDNSAdd.reExecPrivilege(); err != nil {
			return err
		}

		return nil
	}

	// add the DNS entries
	if err := devDNSAdd.addEntries(); err != nil {
		return err
	}

	return nil
}

// validateMeta validates that the required metadata exists
func (devDNSAdd *processDevDNSAdd) validateMeta() error {

	// set the name
	devDNSAdd.name = devDNSAdd.control.Meta["name"]

	// ensure name is provided
	if devDNSAdd.name == "" {
		return fmt.Errorf("Name is required")
	}

	return nil
}

// loadApp loads the app from the db
func (devDNSAdd *processDevDNSAdd) loadApp() error {

	//
	if err := data.Get("apps", config.AppName(), &devDNSAdd.app); err != nil {
		return err
	}

	return nil
}

// entriesExist returns true if both entries already exist
func (devDNSAdd *processDevDNSAdd) entriesExist() bool {

	// generate the entries
	preview := dns.Entry(devDNSAdd.app.DevIP, devDNSAdd.name, "preview")
	dev := dns.Entry(devDNSAdd.app.DevIP, devDNSAdd.name, "dev")

	// if the entries dont exist just return
	if dns.Exists(preview) && dns.Exists(dev) {
		return true
	}

	return false
}

// addEntries adds the dev and preview entries into the host dns
func (devDNSAdd *processDevDNSAdd) addEntries() error {

	// generate the entries
	preview := dns.Entry(devDNSAdd.app.DevIP, devDNSAdd.name, "preview")
	dev := dns.Entry(devDNSAdd.app.DevIP, devDNSAdd.name, "dev")

	// add the 'preview' DNS entry into the /etc/hosts file
	if err := dns.Add(preview); err != nil {
		return err
	}

	// add the 'dev' DNS entry into the /etc/hosts file
	if err := dns.Add(dev); err != nil {
		return err
	}

	return nil
}

// reExecPrivilege re-execs the current process with a privileged user
func (devDNSAdd *processDevDNSAdd) reExecPrivilege() error {

	// call 'dev dns add' with the original path and args; os.Args[0] will be the
	// currently executing program, so this command will ultimately lead right back
	// here
	cmd := fmt.Sprintf("%s dev dns add %s", os.Args[0], devDNSAdd.name)

	// if the sudo'ed subprocess fails, we need to return error to stop the process
	fmt.Println("Admin privileges are required to add DNS entries to your hosts file, your password may be requested...")
	if err := util.PrivilegeExec(cmd); err != nil {
		return err
	}

	return nil
}
