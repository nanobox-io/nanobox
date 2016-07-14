package processor

import ()

// processDestroy ...
type processDestroy struct {
	control ProcessControl
}

//
func init() {
	Register("destroy", destroyFn)
}

//
func destroyFn(control ProcessControl) (Processor, error) {
	return processDestroy{control}, nil
}

//
func (destroy processDestroy) Results() ProcessControl {
	return destroy.control
}

//
func (destroy processDestroy) Process() error {

	// run a provider destroy
	return Run("provider_destroy", destroy.control)
}
