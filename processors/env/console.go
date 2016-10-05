package env

import (
	"fmt"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors/provider"
	"github.com/nanobox-io/nanobox/util/console"
	"github.com/nanobox-io/nanobox/util/display"
)

// Console ...
func Console(componentModel *models.Component, consoleConfig console.ConsoleConfig) error {
	// setup docker client
	if err := provider.Init(); err != nil {
		return err
	}

	// print the MOTD before dropping into the container
	if err := printMOTD(consoleConfig); err != nil {
		return fmt.Errorf("failed to print MOTD: %s", err.Error())
	}

	return console.Run(componentModel.ID, consoleConfig)
}

// printMOTD prints the motd with information for the user to connect
func printMOTD(consoleConfig console.ConsoleConfig) error {

	// print the MOTD
	display.MOTD()

	if consoleConfig.IsDev {
		// print the dev message
		display.InfoDevContainer(consoleConfig.DevIP)
		return nil
	}

	// print the generic message
	display.InfoLocalContainer()
	return nil
}
