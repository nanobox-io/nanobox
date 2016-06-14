// Package dns ...
package dns

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// HOSTSFILE ...
const HOSTSFILE = "/etc/hosts"

// Entry generate the DNS entry to be added
func Entry(ip, name, domain string) string {
	return fmt.Sprintf("%s     %s.%s # '%s' added by running 'nanobox dev dns add <name>'", ip, name, domain, name)
}

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

// Add ...
func Add(entry string) error {

	// open hosts file
	f, err := os.OpenFile(HOSTSFILE, os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	// write the DNS entry to the file
	if _, err := f.WriteString(fmt.Sprintf("%s\n", entry)); err != nil {
		return err
	}

	return nil
}

// Remove ...
func Remove(entry string) error {

	// "contents" will end up storing the entire contents of the file excluding the
	// entry that is trying to be removed
	var contents string

	// open hosts file
	f, err := os.OpenFile(HOSTSFILE, os.O_RDWR, 0644)
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
	if err := ioutil.WriteFile(HOSTSFILE, []byte(contents), 0644); err != nil {
		return err
	}

	return nil
}
