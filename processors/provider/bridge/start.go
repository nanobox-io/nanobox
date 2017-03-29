package bridge

import (
	"github.com/nanobox-io/nanobox/util/provider/bridge"
)

// ask the server to start the bridge
func Start() error {
	return bridge.Start(ConfigFile())
}
