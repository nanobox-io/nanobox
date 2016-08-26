package provider

import (
	"fmt"

	"github.com/jcelliott/lumber"

	"github.com/nanobox-io/nanobox/util/display"
	"github.com/nanobox-io/nanobox/util/locker"
	"github.com/nanobox-io/nanobox/util/provider"
)

// Stop stops the provider (stops the VM)
func Stop() error {
	locker.GlobalLock()
	defer locker.GlobalUnlock()

	display.OpenContext("Halting Nanobox")

	// stop the provider (VM)
	if err := provider.Stop(); err != nil {
		lumber.Error("provider:Stop:provider.Stop(): %s", err.Error())
		return fmt.Errorf("failed to stop the provider: %s", err.Error())
	}

	display.CloseContext()

	return nil
}
