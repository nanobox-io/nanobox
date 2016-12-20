package provider

import (
	"fmt"

	"github.com/jcelliott/lumber"

	"github.com/nanobox-io/nanobox/util/display"
	"github.com/nanobox-io/nanobox/util/locker"
	"github.com/nanobox-io/nanobox/util/provider"
	"github.com/nanobox-io/nanobox/processors/provider/bridge"
)

// Stop stops the provider (stops the VM)
func Stop() error {
	locker.GlobalLock()
	defer locker.GlobalUnlock()

	display.OpenContext("Stopping Nanobox")
	defer display.CloseContext()

	// stop the vpn
	if err := bridge.Stop(); err != nil {
		return fmt.Errorf("failed to stop vpn: %s", err)
	}

	// stop the provider (VM)
	if err := provider.Stop(); err != nil {
		lumber.Error("provider:Stop:provider.Stop(): %s", err.Error())
		return fmt.Errorf("failed to stop the provider: %s", err.Error())
	}

	return nil
}
