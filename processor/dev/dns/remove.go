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

// processDevDNSRemove ...
type processDevDNSRemove struct {
	control processor.ProcessControl
	app     models.App
	name    string
}

func init() {
	processor.Register("dev_dns_remove", devDNSRemoveFn)
}

// devDNSRemveFn creates a processDevDNSRemove and validates the meta in the control
func devDNSRemoveFn(control processor.ProcessControl) (processor.Processor, error) {
	devDNSRemove := &processDevDNSRemove{control: control}
	return devDNSRemove, devDNSRemove.validateMeta()
}

func (devDNSRemove processDevDNSRemove) Results() processor.ProcessControl {
	return devDNSRemove.control
}

//
func (devDNSRemove processDevDNSRemove) Process() error {

	// load the current "app"
	if err := devDNSRemove.loadApp(); err != nil {
		return err
	}

	// short-circuit if the entries dont exist; we do this after we validate meta
	// and load the app because both of those are needed to determin the entry we're
	// looking for
	if !devDNSRemove.entriesExist() {
		return nil
	}

	// if we're not running as the privileged user, we need to re-exec with privilege
	if !util.IsPrivileged() {

		if err := devDNSRemove.reExecPrivilege(); err != nil {
			return err
		}

		return nil
	}

	// remove the DNS entries
	if err := devDNSRemove.removeEntries(); err != nil {
		return err
	}

	return nil
}

// validateMeta validates that the required metadata exists
func (devDNSRemove *processDevDNSRemove) validateMeta() error {

	// set the name
	devDNSRemove.name = devDNSRemove.control.Meta["name"]

	// ensure name is provided
	if devDNSRemove.name == "" {
		return fmt.Errorf("Name is required")
	}

	return nil
}

// loadApp loads the app from the db
func (devDNSRemove *processDevDNSRemove) loadApp() error {

	//
	if err := data.Get("apps", config.AppName(), &devDNSRemove.app); err != nil {
		return err
	}

	return nil
}

// entriesExist returns true if both entries already exist
func (devDNSRemove *processDevDNSRemove) entriesExist() bool {

	// fetch the IPs
	previewIP := devDNSRemove.app.GlobalIPs["preview"]
	devIP := devDNSRemove.app.GlobalIPs["dev"]

	// generate the entries
	preview := dns.Entry(previewIP, devDNSRemove.name, "preview")
	dev := dns.Entry(devIP, devDNSRemove.name, "dev")

	// if the entry doesnt exist just return
	if dns.Exists(preview) && dns.Exists(dev) {
		return true
	}

	return false
}

// removeEntries removes the "dev" and "preview" entries into the host dns
func (devDNSRemove *processDevDNSRemove) removeEntries() error {

	// fetch the IPs
	previewIP := devDNSRemove.app.GlobalIPs["preview"]
	devIP := devDNSRemove.app.GlobalIPs["dev"]

	// generate the entries
	preview := dns.Entry(previewIP, devDNSRemove.name, "preview")
	dev := dns.Entry(devIP, devDNSRemove.name, "dev")

	// remove the 'preview' DNS entry into the /etc/hosts file
	if err := dns.Remove(preview); err != nil {
		return err
	}

	// remove the 'dev' DNS entry into the /etc/hosts file
	if err := dns.Remove(dev); err != nil {
		return err
	}

	return nil
}

// reExecPrivilege re-execs the current process with a privileged user
func (devDNSRemove *processDevDNSRemove) reExecPrivilege() error {

	// call 'dev dns rm' with the original path and args; os.Args[0] will be the
	// currently executing program, so this command will ultimately lead right back
	// here
	cmd := fmt.Sprintf("%s dev dns rm %s", os.Args[0], devDNSRemove.name)

	// if the sudo'ed subprocess fails, we need to return error to stop the process
	fmt.Println("Admin privileges are required to remove DNS entries from your hosts file, your password may be requested...")
	if err := util.PrivilegeExec(cmd); err != nil {
		return err
	}

	return nil
}
