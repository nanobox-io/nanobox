package sim

import (
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors/env"
	"github.com/nanobox-io/nanobox/util/console"
)

// this process is just a shortcut so we can do any other special
// stuff. Which currently there is nothing other then running
// the share console.
func Console(componentModel *models.Component) error {
	return env.Console(componentModel, console.ConsoleConfig{})
}
