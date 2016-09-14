package component

import (
	"fmt"

	"github.com/nanobox-io/nanobox-boxfile"

	"github.com/nanobox-io/nanobox/models"
)

// componentConfig returns the config data from the component boxfile
func componentConfig(component *models.Component) (config map[string]interface{}, err error) {

	// fetch the env
	env, err := models.FindEnvByID(component.EnvID)
	if err != nil {
		err = fmt.Errorf("failed to load env model: %s", err.Error())
		return
	}

	box := boxfile.New([]byte(env.BuiltBoxfile))
	config = box.Node(component.Name).Node("config").Parsed

	switch component.Name {
	case "portal", "logvac", "hoarder", "mist":
		config["token"] = "123"
	}
	return
}
