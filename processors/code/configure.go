package code

import (

	"github.com/jcelliott/lumber"

	"github.com/nanobox-io/nanobox/models"
	code_hook_gen "github.com/nanobox-io/nanobox/generators/hooks/code"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/display"
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

	display.OpenContext("configuring %s", componentModel.Name)
	defer display.CloseContext()

	streamer := display.NewStreamer("info")

	// run fetch build command
	fetchPayload := code_hook_gen.FetchPayload(componentModel, warehouseConfig.WarehouseURL)

	display.StartTask("fetching code")
	if _, err := util.Exec(componentModel.ID, "fetch", fetchPayload, streamer); err != nil {
		display.ErrorTask()
		return err
	}
	display.StopTask()

	// run configure command
	payload := code_hook_gen.ConfigurePayload(appModel, componentModel)

	//
	display.StartTask("configuring code")
	if _, err := util.Exec(componentModel.ID, "configure", payload, streamer); err != nil {
		display.ErrorTask()
		return err
	}
	display.StopTask()

	// run start command
	display.StartTask("starting code")
	if _, err := util.Exec(componentModel.ID, "start", payload, streamer); err != nil {
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
