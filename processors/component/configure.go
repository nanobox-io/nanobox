package component

import (
	"fmt"

	"github.com/jcelliott/lumber"

	"github.com/nanobox-io/nanobox/util/hookit"
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util/display"
	hook_generator "github.com/nanobox-io/nanobox/generators/hooks/component"

)

// Configure configures a component
func Configure(appModel *models.App, componentModel *models.Component) error {
	display.OpenContext("configuring %s(%s)", componentModel.Label, componentModel.Name)
	defer display.CloseContext()

	// short-circuit if the component state is not planned
	if componentModel.State != "planned" {
		return nil
	}

	display.StartTask("running configuration hooks")
	// run the update hook
	if _, err := hookit.RunUpdateHook(componentModel.ID, hook_generator.UpdatePayload(componentModel)); err != nil {
		display.ErrorTask()
		return fmt.Errorf("failed to run update hook: %s", err.Error())
	}

	// run the configure hook
	if _, err := hookit.RunConfigureHook(componentModel.ID, hook_generator.ConfigurePayload(appModel, componentModel)); err != nil {
		display.ErrorTask()
		return fmt.Errorf("failed to run configure hook: %s", err.Error())
	}

	// run the start hook
	if _, err := hookit.RunStartHook(componentModel.ID, hook_generator.UpdatePayload(componentModel)); err != nil {
		display.ErrorTask()
		return fmt.Errorf("failed to run start hook: %s", err.Error())
	}

	// set state as active
	componentModel.State = "active"
	if err := componentModel.Save(); err != nil {
		lumber.Error("component:Configure:models.Component.Save()", err.Error())
		display.ErrorTask()
		return fmt.Errorf("failed to set component state: %s", err.Error())
	}

	display.StopTask()

	return nil
}
