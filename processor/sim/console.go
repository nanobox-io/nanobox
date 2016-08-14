package sim

import (
"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processor/env"
)

// Console ...
type Console struct {
	Component models.Component
}

// this process is just a shortcut so we can do any other special
// stuff. Which currently there is nothing other then running
// the share console.
func (console Console) Run() error {
	envConsole := env.Console{Component: console.Component}
	return envConsole.Run()
}
