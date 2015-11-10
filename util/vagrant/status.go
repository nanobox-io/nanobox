//
package vagrant

import (
	"encoding/json"
	"fmt"
	"github.com/nanobox-io/nanobox/config"
	"io/ioutil"
)

//
type Machine struct {
	Name        string `json:"name"`
	Provider    string `json:"provider"`
	State       string `json:"state"`
	Vagrantfile string `json:"vagrantfile_path"`
}

// Status returns the current status of the VM; this command needs to be run
// in a way independant of a Vagrantfile to ensure that the status will always
// be available
func Status() (status string) {

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
		config.Fatal("[util/vagrant/status] machineIndex failed - ", err.Error())
	}

	// read the machines from machineIndex
	machines := make(map[string]json.RawMessage)
	if err := json.Unmarshal(machineIndex["machines"], &machines); err != nil {
		config.Fatal("[util/vagrant/status] machines failed - ", err.Error())
	}

	// attempt to pull the machine based on the uuid
	machine := Machine{}
	if m, ok := machines[uuid]; ok {
		if err := json.Unmarshal(m, &machine); err != nil {
			config.Fatal("[util/vagrant/status] machine failed - ", err.Error())
		}
	}

	// if the uuid is "" or not found in the machines, then the resulting State
	// will be ""
	status = machine.State

	// if status is "" set it to "not created" to then create the VM
	if status == "" {
		status = "not created"
	}

	return
}
