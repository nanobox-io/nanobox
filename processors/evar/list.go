package evar

import (
	"fmt"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util/odin"
)

func List(envModel *models.Env, appID string) error {
	// fetch the remote
	remote, ok := envModel.Remotes[appID]
	if ok {
		// set the odin endpoint
		odin.SetEndpoint(remote.Endpoint)
		// set the app id
		appID = remote.ID
	}

	evars, err := odin.ListEvars(appID)
	if err != nil {
		return err
	}

	// print the header
	fmt.Printf("\nEnvironment Variables\n")

	// iterate
	for _, evar := range evars {
		fmt.Printf("  %s = %s\n", evar.Key, evar.Value)
	}

	fmt.Println()

	return nil
}
