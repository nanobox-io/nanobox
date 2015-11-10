//
package hosts

import (
	"bufio"
	"fmt"
	"github.com/nanobox-io/nanobox/config"
	"io/ioutil"
	"os"
	"strings"
)

// HasDomain
func HasDomain() (has bool) {

	// open the /etc/hosts file for scanning...
	f, err := os.Open("/etc/hosts")
	if err != nil {
		config.Fatal("[util/file/hosts] os.Open() failed - ", err.Error())
	}
	defer f.Close()

	// scan hosts file looking for an entry corresponding to this app...
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {

		// if an entry with the IP is detected, indicate that it's not needed
		if strings.HasPrefix(scanner.Text(), config.Nanofile.IP) {
			has = true
		}
	}

	return
}

// AddDomain
func AddDomain() {

	// open hosts file
	f, err := os.OpenFile("/etc/hosts", os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		config.Fatal("[util/file/hosts] os.OpenFile() failed - ", err.Error())
	}
	defer f.Close()

	// write the entry to the file
	entry := fmt.Sprintf("\n%-15v   %s # '%v' private network (added by nanobox)", config.Nanofile.IP, config.Nanofile.Domain, config.Nanofile.Name)
	if _, err := f.WriteString(entry); err != nil {
		config.Fatal("[util/file/hosts] file.WriteString() failed - ", err.Error())
	}
}

// RemoveDomain
func RemoveDomain() {

	var contents string

	// open hosts file
	f, err := os.OpenFile("/etc/hosts", os.O_RDWR, 0644)
	if err != nil {
		config.Fatal("[util/file/hosts] os.OpenFile() failed - ", err.Error())
	}
	defer f.Close()

	// remove entry from /etc/hosts
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {

		// if the line contain the entry skip it
		if strings.HasPrefix(scanner.Text(), config.Nanofile.IP) {
			continue
		}

		// add each line back into the file
		contents += fmt.Sprintf("%s\n", scanner.Text())
	}

	// trim the contents to avoid any extra newlines
	contents = strings.TrimSpace(contents)

	// write back the contents of the hosts file minus the removed entry
	if err := ioutil.WriteFile("/etc/hosts", []byte(contents), 0644); err != nil {
		config.Fatal("[util/file/hosts] ioutil.WriteFile failed - ", err.Error())
	}
}
