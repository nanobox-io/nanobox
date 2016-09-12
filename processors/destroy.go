package processors

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jcelliott/lumber"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors/env"
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

	// init docker client
	if err := provider.Init(); err != nil {
		return fmt.Errorf("failed to init docker client: %s", err.Error())
	}

	envModels, _ := models.AllEnvs()
	for _, envModel := range envModels {
		// iterate through the envs and destroy them
		if err := env.Destroy(envModel); err != nil {
			fmt.Printf("unable to remove environment: %s", err)
		}

		// unmount (and remove the share for the env)
		if err := env.Unmount(envModel, false); err != nil {
			fmt.Printf("unable to remove mounts: %s", err)
		}

	}

	// destroy the provider (VM)
	//   this should remove the docker images
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
	if err := clearData(); err != nil {
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

// clearData will remove the global dir and everything inside
func clearData() error {
	dataFile := filepath.ToSlash(filepath.Join(config.GlobalDir(), "data.db"))

	// remove the data.db
	if err := os.Remove(dataFile); err != nil {
		return fmt.Errorf("failed to remove %s: %s", dataFile, err.Error())
	}

	return nil
}
