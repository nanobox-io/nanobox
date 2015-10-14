// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

//
package config

import "os"

// VMfileConfig represents all available/expected .vmfile configurable options
type VMfileConfig struct {
	Deployed    bool   // was the most recent deploy successufl
	Mode        string // foreground/background
	Status      string // the current staus of the VM
	Suspendable bool   // is the VM able to be suspended
}

// ParseVMfile
func ParseVMfile() *VMfileConfig {

	//
	vmfile := &VMfileConfig{}
	vmfilePath := AppDir + "/.vmfile"

	// if a .vmfile doesn't exist - create it
	if _, err := os.Stat(vmfilePath); err != nil {

		vmfile.Deployed = false
		vmfile.Mode = "foreground"
		vmfile.Suspendable = true

		writeVMfile()

		// if a .vmfile does exists - parse it
	} else {
		if err := ParseConfig(vmfilePath, vmfile); err != nil {
			Fatal("[config/vmfile] ParseConfig() failed", err.Error())
		}
	}

	return vmfile
}

//
func (c *VMfileConfig) HasDeployed() bool {
	if err := ParseConfig(AppDir+"/.vmfile", c); err != nil {
		Fatal("[config/vmfile] ParseConfig() failed", err.Error())
	}

	return c.Deployed
}

//
func (c *VMfileConfig) IsMode(mode string) bool {
	if err := ParseConfig(AppDir+"/.vmfile", c); err != nil {
		Fatal("[config/vmfile] ParseConfig() failed", err.Error())
	}

	return c.Mode == mode
}

//
func (c *VMfileConfig) IsSuspendable() bool {
	if err := ParseConfig(AppDir+"/.vmfile", c); err != nil {
		Fatal("[config/vmfile] ParseConfig() failed", err.Error())
	}
	return c.Suspendable
}

//
func (c *VMfileConfig) DeployedIs(deployed bool) {
	c.Deployed = deployed
	writeVMfile()
}

//
func (c *VMfileConfig) ModeIs(mode string) {
	c.Mode = mode
	writeVMfile()
}

//
func (c *VMfileConfig) SuspendableIs(suspendable bool) {
	c.Suspendable = suspendable
	writeVMfile()
}

// writeVMfile
func writeVMfile() {
	if err := writeConfig(AppDir+"/.vmfile", VMfile); err != nil {
		Fatal("[config/vmfile] writeConfig() failed", err.Error())
	}
}
