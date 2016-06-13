package processor

import (
	"github.com/nanobox-io/nanobox/util/counter"
	"github.com/nanobox-io/nanobox/util/data"
)

// processDevReset ...
type processDevReset struct {
	control ProcessControl
}

//
func init() {
	Register("dev_reset", devResetFunc)
}

//
func devResetFunc(control ProcessControl) (Processor, error) {
	return processDevReset{control: control}, nil
}

//
func (devReset processDevReset) Results() ProcessControl {
	return devReset.control
}

//
func (devReset processDevReset) Process() error {

	if err := devReset.resetCounters(); err != nil {
		return err
	}

	return nil
}

// resetCounters resets all the counters associated with all apps
func (devReset processDevReset) resetCounters() error {

	// reset the provider counter
	counter.Reset("provider")

	apps, err := data.Keys("apps")

	if err != nil {
		return err
	}

	for _, app := range apps {
		// reset the general app usage counter
		counter.Reset(app)

		// reset the app dev usage counter
		counter.Reset(app + "_dev")

		// reset the app deploy usage counter
		counter.Reset(app + "_deploy")
	}

	return nil
}
