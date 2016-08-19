package processors

import (
	"fmt"
	
	"github.com/jcelliott/lumber"

	"github.com/nanobox-io/nanobox/processors/provider"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/display"
	"github.com/nanobox-io/nanobox/util/dns"
)

type Destroy struct {}

// Run destroys the provider and cleans nanobox off of the system
func (destroy Destroy) Run() error {
	
	// ensure we're running as the administrator for this
	if !util.IsPrivileged() {
		return destroy.reExecPrivilege()
	}
	
	display.OpenContext("Uninstalling Nanobox and configuration")

	// destroy the provider (VM)
	providerDestroy := provider.Destroy{}
	if err := providerDestroy.Run(); err != nil {
		return fmt.Errorf("failed to uninstall the provider: %s", err.Error())
	}

	display.StartTask("Purging configuration")

	// implode the global dir
	if err := config.ImplodeGlobalDir(); err != nil {
		lumber.Error("Destroy:Run:config.ImplodeGlobalDir(): %s", err.Error())
		return fmt.Errorf("failed to purge the data directory: %s", err.Error())
	}

	// remove all the dns entries
	if err := dns.RemoveAll(); err != nil {
		lumber.Error("Destroy:Run:dns.RemoveAll(): %s", err.Error())
		return fmt.Errorf("failed to remove dns entries: %s", err.Error())
	}
	
	display.StopTask()
	display.CloseContext()
	
	return nil
}

// reExecPrivilege re-execs the current process with a privileged user
func (destroy *Destroy) reExecPrivilege() error {

	display.PrintRequiresPrivilege("to uninstall nanobox and configuration")

	// call this command again, but as superuser
	cmd := fmt.Sprintf("%s destroy", config.NanoboxPath())

	// if the sudo'ed subprocess fails, we need to return error to stop the process
	if err := util.PrivilegeExec(cmd); err != nil {
		lumber.Error("Destroy:reExecPrivilege:util.PrivilegeExec(%s): %s", cmd, err)
		return err
	}
	
	return nil
}
