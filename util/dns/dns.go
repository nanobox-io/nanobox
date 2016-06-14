// Package dns ...
package dns

import (
	"bufio"
	"fmt"
	"os"
)

// HOSTSFILE ...
const HOSTSFILE = "/etc/hosts"

// Exists ...
func Exists(entry string) bool {

	// open the /etc/hosts file for scanning...
	f, err := os.Open(HOSTSFILE)
	if err != nil {
		return false
	}
	defer f.Close()

	// scan each line of the /etc/hosts file to see if there is a match for this
	// entry
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if scanner.Text() == entry {
			return true
		}
	}

	return false
}

// Entry generate the DNS entry to be added
func Entry(ip, name, domain string) string {
	return fmt.Sprintf("%s     %s.%s # '%s' added by running 'nanobox dev dns add <name>'", ip, name, domain, name)
}
