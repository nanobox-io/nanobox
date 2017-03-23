package bridge

import (

	"github.com/nanobox-io/nanobox/commands/server"
)

// ask the server to start the bridge
func Start() error {
	return server.StartBridge(ConfigFile())
}
