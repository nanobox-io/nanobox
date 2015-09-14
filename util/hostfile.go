// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

package util

//
import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/pagodabox/nanobox-cli/config"
	"github.com/pagodabox/nanobox-golang-stylish"
)

// AddDevDomain
func AddDevDomain() {

	// open hosts file
	f, err := os.OpenFile("/etc/hosts", os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		LogFatal("[utils/hostfile] os.OpenFile() failed", err)
	}
	defer f.Close()

	// write the entry to the file
	entry := fmt.Sprintf("\n%-15v   %s # '%v' private network (added by nanobox)", config.Nanofile.IP, config.Nanofile.Domain, config.App)
	if _, err := f.WriteString(entry); err != nil {
		LogFatal("[utils/hostfile] WriteString() failed", err)
	}

	fmt.Printf(stylish.Bullet("Entry for %v (%s.nano.dev) added to /etc/hosts", config.Nanofile.IP, config.App))
}

// RemoveDevDomain
func RemoveDevDomain() {

	// open hosts file
	f, err := os.OpenFile("/etc/hosts", os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		LogFatal("[utils/hostfile] os.OpenFile() failed", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	contents := ""

	// remove entry from /etc/hosts
	for scanner.Scan() {

		// if the line doesn't contain the entry add it back to what is going to be
		// re-written to the file
		if !strings.HasPrefix(scanner.Text(), config.Nanofile.IP) {
			contents += fmt.Sprintf("%s\n", scanner.Text())
		}
	}

	// write back the contents of the hosts file minus the removed entry
	if err := ioutil.WriteFile("/etc/hosts", []byte(contents), 0644); err != nil {
		LogFatal("[utils/hostfile] ioutil.WriteFile failed", err)
	}

	fmt.Printf(stylish.Bullet("Entry for %v (%s.nano.dev) removed from /etc/hosts", config.Nanofile.IP, config.App))
}
