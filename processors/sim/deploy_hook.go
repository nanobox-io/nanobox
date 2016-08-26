package sim

import (
	"encoding/json"

	"github.com/nanobox-io/nanobox-boxfile"
	code_hook_gen "github.com/nanobox-io/nanobox/generators/hooks/code"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util"
)

//
func DeployHook(appModel *models.App, componentModel *models.Component, hooktype string) error {

	_, err := util.Exec(componentModel.ID, hookType, code.DeployPayload(), nil)

	return err
}
