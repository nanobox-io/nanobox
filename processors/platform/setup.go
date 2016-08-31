package platform

import (
	"fmt"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util/display"
)

// Setup provisions platform components needed for an app setup
func Setup(appModel *models.App) error {
	display.OpenContext("Starting platform components")
	defer display.CloseContext()

	for _, component := range setupComponents {
		if err := provisionComponent(appModel, component); err != nil {
			return fmt.Errorf("failed to provision platform component: %s", err.Error())
		}
	}

	return nil
}
