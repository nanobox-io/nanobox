package processors

import (
	"fmt"

	"github.com/nanobox-io/nanobox/helpers"
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/nanoagent"
	"github.com/nanobox-io/nanobox/util/odin"
)

func Tunnel(envModel *models.Env, tunnelConfig TunnelConfig) error {

	appID := tunnelConfig.App

	// fetch the remote
	remote, ok := envModel.Remotes[tunnelConfig.App]
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
		return nil
	}

	// initiate a tunnel session with odin
	key, location, port, err := odin.EstablishTunnel(appID, tunnelConfig.Container)
	if err != nil {
		// todo: can we know if the request was rejected for authorization and print that?
		return fmt.Errorf("failed to initiate a remote tunnel session: %s", err.Error())
	}

	// set a default port if the user didn't specify
	if tunnelConfig.Port == "" {
		tunnelConfig.Port = fmt.Sprintf("%d", port)
	}

	// connect up to the session
	if err := nanoagent.Tunnel(key, location, tunnelConfig.Port); err != nil {
		return fmt.Errorf("failed to connect to remote tunnel session: %s", err.Error())
	}

	return nil
}
