package processor

import "github.com/nanobox-io/nanobox/util/locker"

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

	locker.LocalLock()
	defer locker.LocalUnlock()
	build.control.Meta["build"] = "true"

	// setup the environment (boot vm) but do not run the dev setup because we dont
	// need any o fthe platform services
	if err := Run("provider_setup", build.control); err != nil {
		return err
	}

	// build code
	if err := Run("code_build", build.control); err != nil {
		return err
	}

	return nil
}
