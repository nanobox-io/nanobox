package bridge

import (
	"github.com/nanobox-io/nanobox/commands/server"
)

// ask the server to stop the bridge
func Stop() error {
	return server.StopBridge()
}
