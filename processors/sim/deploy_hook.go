package sim

import (
	generator "github.com/nanobox-io/nanobox/generators/hooks/code"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util/hookit"
)

//
func DeployHook(appModel *models.App, componentModel *models.Component, hookType string) error {

	_, err := hookit.Exec(componentModel.ID, hookType, generator.DeployPayload(appModel, componentModel), "info")

	return err
}
