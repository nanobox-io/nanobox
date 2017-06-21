package processors

import (
	"fmt"
	"os"

	"github.com/nanobox-io/golang-docker-client"

	"github.com/nanobox-io/nanobox/commands/registry"
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors/env"
	"github.com/nanobox-io/nanobox/processors/provider"
	"github.com/nanobox-io/nanobox/processors/server"
	"github.com/nanobox-io/nanobox/util"
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

	// remove any shared caches if any
	docker.VolumeRemove("nanobox_cache")

	// destroy the provider
	if err := provider.Destroy(); err != nil {
		return util.ErrorAppend(err, "failed to destroy the provider")
	}

	// destroy the provider (VM), remove images, remove containers
	if err := util_provider.Implode(); err != nil {
		return util.ErrorAppend(err, "failed to implode the provider")
	}

	// check to see if we need to uninstall nanobox
	// or just remove apps
	if registry.GetBool("full-implode") {

		// teardown the server
		if err := server.Teardown(); err != nil {
			// if we cant tear down the server dont worry about it
			// return util.ErrorAppend(err, "failed to remove server")
		}

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
		return util.ErrorAppend(util.ErrorQuiet(err), "failed to purge the data directory")
	}

	return nil
}
