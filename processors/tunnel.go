package processors

import (
	"fmt"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util/nanoagent"
	"github.com/nanobox-io/nanobox/util/odin"
)

func Tunnel(tunnelConfig TunnelConfig) error {

	// find the app id
	appID := models.AppIDByAlias(tunnelConfig.App)
	if appID == "" {
		// todo: better messaging informing that we couldn't find a link
		return fmt.Errorf("app is not linked")
	}

	// initiate a tunnel session with odin
	key, location, _, port, err := odin.EstablishTunnel(appID, tunnelConfig.Container)
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
