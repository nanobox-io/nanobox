package provider

import (
	"fmt"
	
	"github.com/jcelliott/lumber"
	
	"github.com/nanobox-io/nanobox/util/display"
	"github.com/nanobox-io/nanobox/util/locker"
	"github.com/nanobox-io/nanobox/util/provider"
)

// Destroy destroys the provider
func Destroy() error {
	locker.GlobalLock()
	defer locker.GlobalUnlock()

	display.OpenContext("Destroying Nanobox")

	// destroy the provider
	if err := provider.Destroy(); err != nil {
		lumber.Error("provider:Destroy:provider.Destroy(): %s", err.Error())
		return fmt.Errorf("failed to destroy the provider: %s", err.Error())
	}

	display.CloseContext()

	return nil
}
