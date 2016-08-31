package sim

import (
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors/platform"
)

// Log ...

//
func Log(appModel *models.App) error {
	// some messaging about the logging??
	return platform.MistListen(appModel)
}
