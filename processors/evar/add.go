package evar

import (
	"fmt"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util/display"
	"github.com/nanobox-io/nanobox/util/odin"
)

func Add(envModel *models.Env, appID string, evars map[string]string) error {

	// fetch the link
	link, ok := envModel.Links[appID]
	if ok {
		// set the odin endpoint
		odin.SetEndpoint(link.Endpoint)
		// set the app id
		appID = link.ID
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
