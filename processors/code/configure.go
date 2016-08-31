package code

import (
	"fmt"
	
	"github.com/jcelliott/lumber"

	generator "github.com/nanobox-io/nanobox/generators/hooks/code"
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util/display"
	"github.com/nanobox-io/nanobox/util/hookit"
)

//
func Configure(appModel *models.App, componentModel *models.Component, warehouseConfig WarehouseConfig) error {

	// quit now if the service was activated already
	if componentModel.State == ACTIVE {
		return nil
	}

	// set the prefix so the utilExec lumber logging has context
	lumber.Prefix("code:Configure")
	defer lumber.Prefix("")

	display.OpenContext("configuring %s(%s)", componentModel.Label, componentModel.Name)
	defer display.CloseContext()

	// run fetch build command
	fetchPayload := generator.FetchPayload(componentModel, warehouseConfig.WarehouseURL)

	display.StartTask("fetching code")
	if _, err := hookit.RunFetchHook(componentModel.ID, fetchPayload); err != nil {
		display.ErrorTask()
		return err
	}
	display.StopTask()

	// run configure command
	payload := generator.ConfigurePayload(appModel, componentModel)

	//
	display.StartTask("configuring code")
	if _, err := hookit.RunConfigureHook(componentModel.ID, payload); err != nil {
		display.ErrorTask()
		return fmt.Errorf("failed to configure code: %s", err.Error())
	}
	display.StopTask()

	// run start command
	display.StartTask("starting code")
	if _, err := hookit.RunStartHook(componentModel.ID, payload); err != nil {
		display.ErrorTask()
		return err
	}
	display.StopTask()

	//
	componentModel.State = ACTIVE
	err := componentModel.Save()
	if err != nil {
		lumber.Error("code:Configure:Component.Save(): %s", err.Error())
	}
	return err
}
