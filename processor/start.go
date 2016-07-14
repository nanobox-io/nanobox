package processor

import ()

// processStart ...
type processStart struct {
	control ProcessControl
}

//
func init() {
	Register("start", startFn)
}

//
func startFn(control ProcessControl) (Processor, error) {
	return processStart{control}, nil
}

//
func (start processStart) Results() ProcessControl {
	return start.control
}

//
func (start processStart) Process() error {

	// run a provider setup
	return Run("provider_setup", start.control)
}
