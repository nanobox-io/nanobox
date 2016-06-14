package dns

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

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

	name := devDNSRemove.control.Meta["name"]

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
	if err := devDNSRemove.removeEntry(name, "preview"); err != nil {
		return err
	}

	// remove the 'dev' DNS entry from the /etc/hosts file
	if err := devDNSRemove.removeEntry(name, "dev"); err != nil {
		return err
	}

	return nil
}

// removeEntry ...
func (devDNSRemove processDevDNSRemove) removeEntry(name, domain string) error {

	//
	app := models.App{}
	data.Get("apps", config.AppName(), &app)

	//
	entry := dns.Entry(app.DevIP, name, domain)

	// if the entry doesnt exist just return
	if !dns.Exists(entry) {
		return nil
	}

	// "contents" will end up storing the entire contents of the file excluding the
	// entry that is trying to be removed
	var contents string

	// open hosts file
	f, err := os.OpenFile(dns.HOSTSFILE, os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	// remove entry from /etc/hosts
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {

		// if the line contain the entry skip it
		if scanner.Text() == entry {
			continue
		}

		// add each line back into the file
		contents += fmt.Sprintf("%s\n", scanner.Text())
	}

	// trim the contents to avoid any extra newlines
	contents = strings.TrimSpace(contents)

	// add a single newline for completeness
	contents += "\n"

	// write back the contents of the hosts file minus the removed entry
	if err := ioutil.WriteFile(dns.HOSTSFILE, []byte(contents), 0644); err != nil {
		return err
	}

	return nil
}
