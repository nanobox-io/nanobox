//
package config

import "os"

// VMfileConfig represents all available/expected .vmfile configurable options
type VMfileConfig struct {
	Background  bool // is the CLI running in "background" mode
	Deployed    bool // was the most recent deploy successufl
	Reloaded    bool // did the previous CLI command cause a 'reload'
	Suspendable bool // is the VM able to be suspended
}

// ParseVMfile
func ParseVMfile() (vmfile VMfileConfig) {

	//
	vmfilePath := AppDir + "/.vmfile"

	// if a .vmfile doesn't exist - create it
	if _, err := os.Stat(vmfilePath); err != nil {

		vmfile.Background = false
		vmfile.Deployed = false
		vmfile.Reloaded = false
		vmfile.Suspendable = true

		writeVMfile()
		return
	}

	// if a .vmfile does exists - parse it
	if err := ParseConfig(vmfilePath, &vmfile); err != nil {
		Fatal("[config/vmfile] ParseConfig() failed", err.Error())
	}

	return
}

//
func (c *VMfileConfig) HasDeployed() bool {
	return c.parseVMfile(c.Deployed)
}

//
func (c *VMfileConfig) IsBackground() bool {
	return c.parseVMfile(c.Background)
}

//
func (c *VMfileConfig) HasReloaded() bool {
	return c.parseVMfile(c.Reloaded)
}

//
func (c *VMfileConfig) IsSuspendable() bool {
	return c.parseVMfile(c.Suspendable)
}

//
func (c *VMfileConfig) BackgroundIs(background bool) {
	c.Background = background
	writeVMfile()
}

//
func (c *VMfileConfig) DeployedIs(deployed bool) {
	c.Deployed = deployed
	writeVMfile()
}

//
func (c *VMfileConfig) ReloadedIs(reloaded bool) {
	c.Reloaded = reloaded
	writeVMfile()
}

//
func (c *VMfileConfig) SuspendableIs(suspendable bool) {
	c.Suspendable = suspendable
	writeVMfile()
}

// parseVMfile is a wrapper that simply handles the error once rather than in
// each individual call
func (c *VMfileConfig) parseVMfile(field bool) bool {
	if err := ParseConfig(AppDir+"/.vmfile", c); err != nil {
		Fatal("[config/vmfile] ParseConfig() failed", err.Error())
	}

	return field
}

// writeVMfile writes to the vmfile with each field update
func writeVMfile() {
	if err := writeConfig(AppDir+"/.vmfile", VMfile); err != nil {
		Fatal("[config/vmfile] writeConfig() failed", err.Error())
	}
}
