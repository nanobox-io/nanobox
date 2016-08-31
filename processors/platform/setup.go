package platform

import (
	"fmt"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util/display"
)

// Setup provisions platform components needed for an app setup
func Setup(a *models.App) error {
	display.OpenContext("provisioning platform services")
	defer display.CloseContext()

	for _, component := range setupComponents {
		if err := provisionComponent(a, component); err != nil {
			return fmt.Errorf("failed to provision platform component: %s", err.Error())
		}
	}

	return nil
}
