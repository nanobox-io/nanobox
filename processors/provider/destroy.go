package provider

import (
	"github.com/jcelliott/lumber"

	"github.com/nanobox-io/nanobox/processors/provider/bridge"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/locker"
	"github.com/nanobox-io/nanobox/util/provider"
)

// Destroy destroys the provider
func Destroy() error {
	locker.GlobalLock()
	defer locker.GlobalUnlock()

	if provider.BridgeRequired() {

		// remove the network bridge
		if err := bridge.Teardown(); err != nil {
			return util.ErrorAppend(err, "failed to teardown network bridge")
		}
	}

	// destroy the provider
	if err := provider.Destroy(); err != nil {
		lumber.Error("provider:Destroy:provider.Destroy()")
		return util.ErrorAppend(err, "failed to destroy the provider")
	}

	return nil
}
