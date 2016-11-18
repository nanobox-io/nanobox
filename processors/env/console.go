package env

import (
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

	switch {
	case consoleConfig.Command != "":
		display.InfoDevRunContainer(consoleConfig.Command, consoleConfig.DevIP)
	default:
		display.MOTD()
		display.InfoDevContainer(consoleConfig.DevIP)
	}

	return console.Run(componentModel.ID, consoleConfig)
}
