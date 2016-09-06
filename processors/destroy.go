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

// Destroy destroys the provider and cleans nanobox off of the system
func Destroy() error {

	// ensure we're running as the administrator for this
	if !util.IsPrivileged() {
		return reExecPrivilegedDestroy()
	}

	display.OpenContext("Uninstalling Nanobox")
	defer display.CloseContext()

	// todo: we need to make this process a bit more robust once a VM isn't
	// our only provider. We won't be able to just "destroy" the provider once
	// we're dealing with native providers. The following todos will address
	// this scenario:

	// todo: iterate through the envs and destroy them

	// todo: unmount (and remove the share for the env)

	// todo: purge the installed docker images

	// destroy the provider (VM)
	if err := provider.Destroy(); err != nil {
		return fmt.Errorf("failed to uninstall the provider: %s", err.Error())
	}

	// purge the installation
	if err := purgeConfiguration(); err != nil {
		return fmt.Errorf("failed to purge nanobox configuration: %s", err.Error())
	}

	return nil
}

// purges the config data and dns entries
func purgeConfiguration() error {
	display.StartTask("Purging configuration")
	defer display.StopTask()

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

	return nil
}

// reExecPrivilege re-execs the current process with a privileged user
func reExecPrivilegedDestroy() error {
	display.PauseTask()
	defer display.ResumeTask()

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
