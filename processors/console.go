package processors

import (
	"fmt"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util/nanoagent"
	"github.com/nanobox-io/nanobox/util/odin"
)

func Console(app, host string) error {

	// lookup the app by the alias
	appID := models.AppIDByAlias(app)
	if appID == "" {
		// todo: better messaging informing that we couldn't find a link
		return fmt.Errorf("app is not linked")
	}

	// initiate a console session with odin
	key, location, container, err := odin.EstablishConsole(appID, host)
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
