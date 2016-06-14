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

// processDevDNSRemove
type processDevDNSRemove struct {
	control processor.ProcessControl
}

//
func init() {
	processor.Register("dev_dns_remove", devDNSRemoveFunc)
}

//
func devDNSRemoveFunc(control processor.ProcessControl) (processor.Processor, error) {
	return processDevDNSRemove{control: control}, nil
}

//
func (devDNSRemove processDevDNSRemove) Results() processor.ProcessControl {
	return devDNSRemove.control
}

//
func (devDNSRemove processDevDNSRemove) Process() error {

	//
	name := devDNSRemove.control.Meta["name"]

	//
	app := models.App{}
	data.Get("apps", config.AppName(), &app)

	//
	entry := dns.Entry(app.DevIP, name, domain)

	// if the entry already exists just return
	if dns.Exists(entry) {
		return nil
	}

	// This process requires root, check to see if we're the root user. If not, we
	// need to run a hidden command as sudo that will just call this function again.
	// Thus, the subprocess will be running as root
	if os.Geteuid() != 0 {

		// get the original nanobox executable
		nanobox := os.Args[0]

		// call dev netfs add with the original path (ultimately leads right back here)
		cmd := fmt.Sprintf("%s dev dns rm %s", nanobox, name)

		// if the sudo'ed subprocess fails, we need to return error to stop the process
		fmt.Println("Admin privileges are required to remove DNS entries from your hosts file, your password may be requested...")
		if err := util.PrivilegeExec(cmd); err != nil {
			return err
		}

		// the subprocess exited successfully, so we can short-circuit here
		return nil
	}

	// remove the 'preview' DNS entry from the /etc/hosts file
	if err := dns.Remove(name, "preview"); err != nil {
		return err
	}

	// remove the 'dev' DNS entry from the /etc/hosts file
	if err := dns.Remove(name, "dev"); err != nil {
		return err
	}

	return nil
}
