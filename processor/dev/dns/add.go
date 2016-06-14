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
	app			models.App
}

//
func init() {
	processor.Register("dev_dns_add", devDNSAddFunc)
}

//
func devDNSAddFunc(control processor.ProcessControl) (processor.Processor, error) {
	return processDevDNSAdd{control: control}, nil
}

//
func (devDNSAdd processDevDNSAdd) Results() processor.ProcessControl {
	return devDNSAdd.control
}

//
func (devDNSAdd processDevDNSAdd) Process() error {

	if err := devDNSAdd.validateMeta(); err != nil {
		return err
	}

	if err := devDNSAdd.loadApp(); err != nil {
		return err
	}

	// short-circuit if the entries already exist
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

	if err := devDNSAdd.addEntries(); err != nil {
		return err
	}

	return nil
}

// validateMeta validates that the required metadata exists
func (devDNSAdd processDevDNSAdd) validateMeta() error {

	// ensure name is provided
	if devDNSAdd.control.Meta["name"] == "" {
		return errors.New("Name is required")
	}

	return nil
}

// loadApp loads the app from the db
func (devDNSAdd processDevDNSAdd) loadApp() error {

	if err := data.Get("apps", config.AppName(), &devDNSAdd.app); err != nil {
		return err
	}

	return nil
}

// entriesExist returns true if both entries already exist
func (devDNSAdd processDevDNSAdd) entriesExist() bool {

	name := devDNSAdd.control.Meta["name"]

	// generate the entries
	preview := dns.Entry(devDNSAdd.app.DevIP, name, "preview")
	dev := dns.Entry(devDNSAdd.app.DevIP, name, "dev")

	// if the entry doesnt exist just return
	if dns.Exists(preview) && dns.Exists(dev) {
		return true
	}

	return false
}

// addEntries adds the dev and preview entries into the host dns
func (devDNSAdd processDevDNSAdd) addEntries() error {
	name := devDNSAdd.control.Meta["name"]

	// generate the entries
	preview := dns.Entry(devDNSAdd.app.DevIP, name, "preview")
	dev := dns.Entry(devDNSAdd.app.DevIP, name, "dev")

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
func (devDNSAdd processDevDNSAdd) reExecPrivilege() error {
	name := devDNSAdd.control.Meta["name"]

	// get the original nanobox executable
	nanobox := os.Args[0]

	// call 'dev dns add' with the original path (ultimately leads right back here)
	cmd := fmt.Sprintf("%s dev dns add %s", nanobox, name)

	// if the sudo'ed subprocess fails, we need to return error to stop the process
	fmt.Println("Admin privileges are required to add DNS entries to your hosts file, your password may be requested...")
	if err := util.PrivilegeExec(cmd); err != nil {
		return err
	}

	return nil
}
