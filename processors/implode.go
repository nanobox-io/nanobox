package processors

import (
	"fmt"
	"os"

	"github.com/jcelliott/lumber"

	"github.com/nanobox-io/nanobox/commands/registry"
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors/env"
	"github.com/nanobox-io/nanobox/processors/provider"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/display"
	util_provider "github.com/nanobox-io/nanobox/util/provider"
)

// Implode destroys the provider and cleans nanobox off of the system
func Implode() error {

	display.OpenContext("Imploding Nanobox")
	defer display.CloseContext()

	// remove all environments
	envModels, _ := models.AllEnvs()
	for _, envModel := range envModels {
		// remove all environments
		if err := env.Destroy(envModel); err != nil {
			fmt.Printf("unable to remove mounts: %s", err)
		}
	}

	// destroy the provider
	if err := provider.Destroy(); err != nil {
		return fmt.Errorf("failed to destroy the provider: %s", err)
	}

	// destroy the provider (VM), remove images, remove containers
	if err := util_provider.Implode(); err != nil {
		return fmt.Errorf("failed to implode the provider: %s", err)
	}

	// purge the installation
	if registry.GetBool("full-implode") {
		purgeConfiguration()
	}

	return nil
}

// purges the config data and dns entries
func purgeConfiguration() error {

	display.StartTask("Purging configuration")
	defer display.StopTask()

	// implode the global dir
	if err := os.RemoveAll(config.GlobalDir()); err != nil {
		lumber.Error("Destroy:Run:config.ImplodeGlobalDir(): %s", err.Error())
		return fmt.Errorf("failed to purge the data directory: %s", err.Error())
	}

	return nil
}
