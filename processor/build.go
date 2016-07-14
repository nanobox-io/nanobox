package processor

import (
	"github.com/nanobox-io/nanobox/util/locker"
)

// processBuild ...
type processBuild struct {
	control ProcessControl
}

//
func init() {
	Register("build", buildFn)
}

//
func buildFn(control ProcessControl) (Processor, error) {
	return processBuild{control}, nil
}

//
func (build processBuild) Results() ProcessControl {
	return build.control
}

//
func (build processBuild) Process() error {

	// by aquiring a local lock we are only allowing
	// one build to happen at a time
	locker.LocalLock()
	defer locker.LocalUnlock()

	// get the vm and app up.
	if err := Run("env_setup", build.control); err != nil {
		return err
	}

	// build code
	if err := Run("code_build", build.control); err != nil {
		return err
	}

	return nil
}
