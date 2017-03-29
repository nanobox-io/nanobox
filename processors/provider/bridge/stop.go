package bridge

import (
	"github.com/nanobox-io/nanobox/util/provider/bridge"
)

// ask the server to stop the bridge
func Stop() error {
	return bridge.Stop()
}
