// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

//
package vagrant

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/nanobox-io/nanobox-cli/config"
)

// Status returns the current status of the VM; this command needs to be run
// in a way independant of a Vagrantfile to ensure that the status will always
// be available
func Status() string {

	// attempt to get the uuid; don't really care about the error here since the
	// value will be "" if there is no file to read
	b, _ := ioutil.ReadFile(fmt.Sprintf("%v/.vagrant/machines/%v/%v/index_uuid", config.AppDir, config.Nanofile.Name, config.Nanofile.Provider))

	// set uuid
	uuid := string(b)

	// attempt to get the machine data; don't really care about the error here since
	// the value will be "" if there is no file to read
	b, _ = ioutil.ReadFile(config.Home + "/.vagrant.d/data/machine-index/index")

	// get the initial data
	machineIndex := make(map[string]json.RawMessage)
	if err := json.Unmarshal(b, &machineIndex); err != nil {
		config.Fatal("[util/vagrant/status] failed - ", err.Error())
	}

	// read the machines from machineIndex
	machines := make(map[string]json.RawMessage)
	if err := json.Unmarshal(machineIndex["machines"], &machines); err != nil {
		config.Fatal("[util/vagrant/status] failed - ", err.Error())
	}

	// attempt to pull the machine based on the uuid
	machine := Machine{}
	if err := json.Unmarshal(machines[uuid], &machine); err != nil {
		config.Fatal("[util/vagrant/status] failed - ", err.Error())
	}

	// if the uuid is "" or not found in the machines, then the resulting State
	// will be ""
	status := machine.State

	// set the status in the .vmfile (currently not used)
	// config.VMfile.StatusIs(status)

	if machine.State == "" {
		status = "not created"
	}

	return status
}
