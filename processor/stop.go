package processor

import (
	"github.com/nanobox-io/nanobox/processor/app"
	"github.com/nanobox-io/nanobox/processor/provider"
)

// Stop ...
type Stop struct {
}

//
func (stop Stop) Run() error {

	// stop all running environments
	if err := stop.stopAllApps(); err != nil {
		return err
	}

	// run a provider setup
	providerStop := provider.Stop{}
	return providerStop.Run()
}

// stop all of the apps that are currently up
func (stop Stop) stopAllApps() error {

	// run the app stop on all running apps
	for _, a := range upApps() {

		appStop := app.Stop{
			App: a,
		}
		
		if err := appStop.Run(); err != nil {
			return err
		}

	}
	return nil
}
