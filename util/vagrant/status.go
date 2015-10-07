// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

//
package vagrant

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os/exec"
	"regexp"
	"strings"

	"github.com/nanobox-io/nanobox-cli/config"
)

// Status returns the current status of the VM; this command needs to be run
// in a way independant of a Vagrantfile to ensure that the status will always
// be available
func Status() string {

	var status string

	// don't really care about the error here since the value will be "" if there
	// is no file to read
	b, _ := ioutil.ReadFile(fmt.Sprintf("%v/.vagrant/machines/%v/%v/index_uuid", config.AppDir, config.Nanofile.Name, config.Nanofile.Provider))
	uuid := string(b)

	// if no uuid is found the vm has not yet been created
	if uuid == "" {
		return "not created"
	}

	// run the command with the uuid; this allows the command to be run from any
	// context (not dependent on a Vagrantfile)
	cmd := exec.Command("vagrant", "status", uuid)

	//
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		config.Fatal("[util/vagrant] cmd.StdoutPipe() failed", err.Error())
	}

	//
	scanner := bufio.NewScanner(stdout)
	go func() {
		for scanner.Scan() {

			// if the scan line starts with the name of the VM thats the line thats going
			// to have the status in it
			if strings.HasPrefix(scanner.Text(), config.Nanofile.Name) {

				// pull the status out of the line; it's easiest to just pull the exact
				// states out of the string since there aren't that many
				submatch := regexp.MustCompile(`\s(running|restoring|saved|poweroff|not created|aborted)\s`).FindStringSubmatch(scanner.Text())

				//
				if len(submatch) > 1 {
					status = submatch[1]
				} else {
					config.Fatal("[util/vagrant] Unknown status: ", scanner.Text())
				}
			}
		}
	}()

	// start the command; we need this to 'fire and forget' so that we can manually
	// capture and modify the commands output
	if err := cmd.Start(); err != nil {
		config.Fatal("[util/vagrant] cmd.Start() failed", err.Error())
	}

	// wait for the command to complete/fail and exit, giving us a chance to scan
	// the output
	if err := cmd.Wait(); err != nil {
		config.Fatal("[util/vagrant] cmd.Wait() failed", err.Error())
	}

	// set the status in the .vmfile (currently not used)
	// config.VMfile.StatusIs(status)

	return status
}
