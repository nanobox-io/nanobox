package processor

import ()

// processStop ...
type processStop struct {
	control ProcessControl
}

//
func init() {
	Register("stop", stopFn)
}

//
func stopFn(control ProcessControl) (Processor, error) {
	return processStop{control}, nil
}

//
func (stop processStop) Results() ProcessControl {
	return stop.control
}

//
func (stop processStop) Process() error {

	// stop all running environments
	if err := stop.stopAllApps(); err != nil {
		return err
	}

	// run a provider setup
	return Run("provider_stop", stop.control)
}

// stop all of the apps that are currently up
func (stop processStop) stopAllApps() error {

	// create a control for the child processes
	control := stop.control.Dup()

	// run the app stop on all running apps
	for _, app := range upApps() {
		control.Meta["name"] = app.ID

		err := Run("app_stop", control)
		if err != nil {
			return err
		}

	}
	return nil
}
