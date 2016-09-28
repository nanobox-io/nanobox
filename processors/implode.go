package processors

import (
	"fmt"
	"os"

	"github.com/jcelliott/lumber"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors/env"
	"github.com/nanobox-io/nanobox/processors/provider"
	util_provider "github.com/nanobox-io/nanobox/util/provider"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/display"
	"github.com/nanobox-io/nanobox/util/dns"
)

// Implode destroys the provider and cleans nanobox off of the system
func Implode() error {

	display.OpenContext("Imploding Nanobox")
	defer display.CloseContext()

	// init docker client
	if err := provider.Init(); err != nil {
		return fmt.Errorf("failed to init docker client: %s", err.Error())
	}

	envModels, _ := models.AllEnvs()
	for _, envModel := range envModels {
		// unmount (and remove the share for the env)
		if err := env.Unmount(envModel, false); err != nil {
			fmt.Printf("unable to remove mounts: %s", err)
		}

	}

	// destroy the provider (VM), remove images, remove containers
	if err := util_provider.Implode(); err != nil {
		return fmt.Errorf("failed to implode the provider: %s", err.Error())
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

// clearData will remove the global dir and everything inside
func clearData() error {

	// remove the .nanobox/ folder
	if err := os.RemoveAll(config.GlobalDir()); err != nil {
		return fmt.Errorf("failed to remove %s: %s", config.GlobalDir(), err.Error())
	}

	return nil
}
