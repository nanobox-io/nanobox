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

	// fetch the link
	link, ok := envModel.Links[consoleConfig.App]
	if ok {
		// set the odin endpoint
		odin.SetEndpoint(link.Endpoint)
		// set the app id
		appID = link.ID
	}

	// if an endpoint was provided as a flag, override the linked endpoint
	if consoleConfig.Endpoint != "" {
		odin.SetEndpoint(consoleConfig.Endpoint)
	}
	
	// set the app id to the directory name if it's default
	if appID == "default" {
		appID = config.AppName()
	}
	
	// validate access to the app
	if err := helpers.ValidateOdinApp(appID); err != nil {
		// the validation already printed the error
		return nil
	}

	// initiate a console session with odin
	key, location, container, err := odin.EstablishConsole(appID, consoleConfig.Host)
	if err != nil {
		// todo: can we know if the request was rejected for authorization and print that?
		return fmt.Errorf("failed to initiate a remote console session: %s", err.Error())
	}

	// todo: extract the protocol from above and run ssh or console

	// connect up to the session
	if err = nanoagent.Console(key, location, container); err != nil {
		return fmt.Errorf("failed to connect to remote console session: %s", err.Error())
	}

	return nil
}
