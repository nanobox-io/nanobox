package evar

import (
	"fmt"

	"github.com/nanobox-io/nanobox/commands/registry"
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util/display"
	"github.com/nanobox-io/nanobox/util/odin"
)

func Add(envModel *models.Env, appID string, evars map[string]string) error {

	// fetch the remote
	remote, ok := envModel.Remotes[appID]
	if ok {
		// set the odin endpoint
		odin.SetEndpoint(remote.Endpoint)
		// set the app id
		appID = remote.ID
	}

	// set odins endpoint if the arguement is passed
	if endpoint := registry.GetString("endpoint"); endpoint != "" {
		odin.SetEndpoint(endpoint)
	}

	// iterate through the evars and add them to the app
	for key, val := range evars {
		err := odin.AddEvar(appID, key, val)
		if err != nil {
			return err
		}
		fmt.Printf("%s %s added\n", display.TaskComplete, key)
	}

	return nil
}
