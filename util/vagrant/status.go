//
package vagrant

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/nanobox-io/nanobox/config"
)

// Status returns the current status of the VM; this command needs to be run
// in a way independant of a Vagrantfile to ensure that the status will always
// be available
func Status() (status string) {

	// default status of the VM
	status = "not created"

	// attempt to get the uuid; don't handle the error here because there are some
	// other conditions we want lumped together
	b, _ := ioutil.ReadFile(fmt.Sprintf("%v/.vagrant/machines/%v/%v/index_uuid", config.AppDir, config.Nanofile.Name, config.Nanofile.Provider))

	// set uuid (this will be "" if the above returned an error)
	uuid := string(b)

	// attempt to get the machine data
	b, err = ioutil.ReadFile(config.Home + "/.vagrant.d/data/machine-index/index")

	// an error here (os.PathError) means that the files was not found, causing
	// meaning no VMs have been created. Returning here will indicate a VM needs
	// to be created due to the default status above
	if uuid == "" || err != nil {
		return
	}

	// get the initial data
	machineIndex := make(map[string]json.RawMessage)
	if err := json.Unmarshal(b, &machineIndex); err != nil {
		config.Fatal("[util/vagrant/status] json.Unmarshal() machineIndex failed", err.Error())
	}

	// read the machines from machineIndex
	machines := make(map[string]json.RawMessage)
	if err := json.Unmarshal(machineIndex["machines"], &machines); err != nil {
		config.Fatal("[util/vagrant/status] json.Unmarshal() machines failed", err.Error())
	}

	// attempt to pull the machine based on the uuid
	machine := struct {
		Name        string `json:"name"`
		Provider    string `json:"provider"`
		State       string `json:"state"`
		Vagrantfile string `json:"vagrantfile_path"`
	}{}
	if m, ok := machines[uuid]; ok {
		if err := json.Unmarshal(m, &machine); err != nil {
			config.Fatal("[util/vagrant/status] json.Unmarshal() machine failed", err.Error())
		}
	}

	//
	if machine.State != "" {
		status = machine.State
	}

	return
}
