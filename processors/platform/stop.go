package platform

import (
	"github.com/jcelliott/lumber"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors/component"
	"github.com/nanobox-io/nanobox/util"
)

// Stop stops all platform components
func Stop(a *models.App) error {
	for _, pc := range setupComponents {
		if err := stopComponent(a, pc); err != nil {
			return util.ErrorAppend(err, "failed to stop platform component")
		}
	}

	return nil
}

// stopComponent stops a platform component
func stopComponent(a *models.App, pc PlatformComponent) error {
	// load the component
	c, err := models.FindComponentBySlug(a.ID, pc.name)
	if err != nil {
		lumber.Error("platform:stopComponent:models.FindComponentBySlug(%s, %s): %s", a.ID, pc.name, err.Error())
		return util.ErrorAppend(err, "failed to load component")
	}

	// stop the component
	if err := component.Stop(c); err != nil {
		return util.ErrorAppend(err, "failed to stop component")
	}

	return nil
}
