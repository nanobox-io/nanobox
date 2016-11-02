package evar

import (
	"fmt"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util/display"
	"github.com/nanobox-io/nanobox/util/odin"
)

func Remove(envModel *models.Env, appID string, keys []string) error {

	// fetch the link
	link, ok := envModel.Links[appID]
	if ok {
		// set the odin endpoint
		odin.SetEndpoint(link.Endpoint)
		// set the app id
		appID = link.ID
	}

	evars, err := odin.ListEvars(appID)
	if err != nil {
		return err
	}

	// delete the evars
	for _, key := range keys {
		removed := false
		for _, evar := range evars {
			if evar.Key == key {
				if err := odin.RemoveEvar(appID, evar.ID); err != nil {
					return err
				}
				removed = true
				fmt.Printf("%s %s removed\n", display.TaskComplete, key)
			}
		}
		if !removed {
			fmt.Printf("%s %s not found\n", display.TaskPause, key)
		}
	}
	fmt.Println()

	return nil
}
