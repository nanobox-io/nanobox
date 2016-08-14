package provider

import (
	"github.com/nanobox-io/nanobox/provider"
	"github.com/nanobox-io/nanobox/util/locker"
)

// Destroy ...
type Destroy struct {
}

//
func (destroy Destroy) Run() error {
	locker.GlobalLock()
	defer locker.GlobalUnlock()

	//
	if err := destroy.removeDatabase(); err != nil {
		return err
	}

	return provider.Destroy()
}
