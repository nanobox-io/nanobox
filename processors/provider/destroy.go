package provider

import (
	"github.com/nanobox-io/nanobox/util/provider"
	"github.com/nanobox-io/nanobox/util/display"
	"github.com/nanobox-io/nanobox/util/locker"
)

// Destroy ...
type Destroy struct {
}

//
func (destroy Destroy) Run() error {
	locker.GlobalLock()
	defer locker.GlobalUnlock()

	// delete the nanobox database
	display.StartTask("removing database")
	if err := destroy.removeDatabase(); err != nil {
		return err
	}
	display.StopTask()

	// remove the provider
	display.StartTask("removing vm")
	err := provider.Destroy()
	if err != nil {
		display.ErrorTask()
		return err
	}
	display.StopTask()

	return nil
}

