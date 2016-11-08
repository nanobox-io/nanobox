package processors

import (
	"fmt"

	"github.com/nanobox-io/nanobox/helpers"
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/nanoagent"
	"github.com/nanobox-io/nanobox/util/odin"
)

func Console(envModel *models.Env, consoleConfig ConsoleConfig) error {

	appID := consoleConfig.App

	// fetch the remote
	remote, ok := envModel.Remotes[appID]
	if ok {
		// set the odin endpoint
		odin.SetEndpoint(remote.Endpoint)
		// set the app id
		appID = remote.ID
	}
	
	// set the app id to the directory name if it's default
	if appID == "default" {
		appID = config.AppName()
	}

	// validate access to the app
	if err := helpers.ValidateOdinApp(appID); err != nil {
		// the validation already printed the error
		return err
	}

	// initiate a console session with odin
	key, location, protocol, err := odin.EstablishConsole(appID, consoleConfig.Host)
	if err != nil {
		// todo: can we know if the request was rejected for authorization and print that?
		return fmt.Errorf("failed to initiate a remote console session: %s", err.Error())
	}

	switch protocol {
	case "docker":
		if err := nanoagent.Console(key, location); err != nil {
			return fmt.Errorf("failed to connect to remote console session: %s", err.Error())
		}
	case "ssh":
		if err := nanoagent.SSH(key, location); err != nil {
			return fmt.Errorf("failed to connect to remote ssh server: %s", err)
		}
	}

	return nil
}
