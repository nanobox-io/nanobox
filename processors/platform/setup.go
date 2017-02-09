package platform

import (
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/display"
)

// Setup provisions platform components needed for an app setup
func Setup(appModel *models.App) error {
	display.OpenContext("Starting components")
	defer display.CloseContext()

	for _, component := range setupComponents {
		if err := provisionComponent(appModel, component); err != nil {
			return util.ErrorAppend(err, "failed to provision platform component")
		}
	}

	return nil
}
