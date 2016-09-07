package processors

import (
	"fmt"

	"github.com/nanobox-io/nanobox/helpers"
	"github.com/nanobox-io/nanobox/util/nanoagent"
	"github.com/nanobox-io/nanobox/util/odin"
)

func Tunnel(tunnelConfig TunnelConfig) error {

	appID, err := helpers.OdinAppIDByAlias(tunnelConfig.App)
	if err != nil {
		// the message will have already been printed in the helper
		return nil
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
