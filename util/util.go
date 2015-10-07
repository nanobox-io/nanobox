// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

//
package util

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/nanobox-io/nanobox-cli/config"
)

// NeedsDomain
func NeedsDomain() (need bool) {

	// open the /etc/hosts file for scanning...
	f, err := os.Open("/etc/hosts")
	if err != nil {
		config.Fatal("[commands/create] os.Open() failed", err.Error())
	}
	defer f.Close()

	// determines whether or not an entry needs to be added to the /etc/hosts file
	// (an entry will be added unless it's confirmed that it's not needed)
	need = true

	// scan hosts file looking for an entry corresponding to this app...
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {

		// if an entry with the IP is detected, flag the entry as not needed
		if strings.HasPrefix(scanner.Text(), config.Nanofile.IP) {
			need = false
		}
	}

	return
}

// HostfileAddDomain
func HostfileAddDomain() {

	// open hosts file
	f, err := os.OpenFile("/etc/hosts", os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		config.Fatal("[utils/hostfile] os.OpenFile() failed", err.Error())
	}
	defer f.Close()

	// write the entry to the file
	entry := fmt.Sprintf("\n%-15v   %s # '%v' private network (added by nanobox)", config.Nanofile.IP, config.Nanofile.Domain, config.Nanofile.Name)
	if _, err := f.WriteString(entry); err != nil {
		config.Fatal("[utils/hostfile] WriteString() failed", err.Error())
	}
}

// HostfileRemoveDomain
func HostfileRemoveDomain() {

	var contents string

	// open hosts file
	f, err := os.OpenFile("/etc/hosts", os.O_RDWR, 0644)
	if err != nil {
		config.Fatal("[utils/hostfile] os.OpenFile() failed", err.Error())
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
		config.Fatal("[utils/hostfile] ioutil.WriteFile failed", err.Error())
	}
}

//
func MD5sMatch(localPath, remotePath string) (bool, error) {

	// get local md5
	f, err := os.Open(localPath)

	// if there is no local md5 return false
	if err != nil {
		return false, nil
	}
	defer f.Close()

	localMD5, err := ioutil.ReadAll(f)
	if err != nil {
		return false, err
	}

	// get remote md5
	res, err := http.Get(remotePath)
	if err != nil {
		return false, err
	}
	defer res.Body.Close()

	remoteMD5, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return false, err
	}

	return string(localMD5) == string(remoteMD5), nil
}
