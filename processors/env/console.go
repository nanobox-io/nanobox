package env

import (
	"time"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors/provider"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/console"
	"github.com/nanobox-io/nanobox/util/display"
)

// Console ...
func Console(componentModel *models.Component, consoleConfig console.ConsoleConfig) error {
	if componentModel.ID == "" {
		display.ConsoleNodeNotFound()
		return util.Err{
			Message: "Node not found",
			Code:    "USER",
			Stack:   []string{"failed to console"},
			Suggest: "It appears the node specified does not exist. Please double check the node name in your boxfile.yml.",
		}
	}
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
	<-time.After(100 * time.Millisecond)

	return console.Run(componentModel.ID, consoleConfig)
}
