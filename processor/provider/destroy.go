package provider

import (
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processor/app/dns"
	"github.com/nanobox-io/nanobox/provider"
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

	// clean all the entries in the /etc/hosts file
	// for all apps
	display.StartTask("cleaning hosts")
	removeAllDns := dns.RemoveAll{
		// by setting the id to "by nanobox" removeAll
		// will remove all entries that were set by us
		App: models.App{ID: "by nanobox"},
	}
	if err := removeAllDns.Run(); err != nil {
		display.ErrorTask()
		return err
	}
	display.StopTask()

	return nil
}
