package platform

import (
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors/component"
	"github.com/nanobox-io/nanobox/util"
)

// provisionComponent will provision an individual component
func provisionComponent(appModel *models.App, platformComponent PlatformComponent) error {

	componentModel := &models.Component{
		Name:  platformComponent.name,
		Label: platformComponent.label,
		Image: platformComponent.image,
	}

	// if the component exists and is active just start it and return
	if isComponentActive(appModel, componentModel) {

		// start the component
		if err := component.Start(componentModel); err != nil {
			return util.ErrorAppend(err, "failed to start component")
		}

		return nil
	}

	// setup
	if err := component.Setup(appModel, componentModel); err != nil {
		return util.ErrorAppend(err, "failed to setup platform component (%s)", componentModel.Label)
	}

	return nil
}

// isComponentActive returns true if a component is already active
func isComponentActive(appModel *models.App, componentModel *models.Component) bool {
	// component db entry
	component, _ := models.FindComponentBySlug(appModel.ID, componentModel.Name)
	if component.State == "active" {

		// set the componentModel pointer to the new component object
		*componentModel = *component
		return true
	}

	return false
}
