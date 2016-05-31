package processor

import (
	"github.com/nanobox-io/nanobox/util/locker"
)

type build struct {
	control ProcessControl
}

func init() {
	Register("build", buildFunc)
}

func buildFunc(control ProcessControl) (Processor, error) {
	return build{control}, nil
}

func (self build) Results() ProcessControl {
	return self.control
}

func (self build) Process() error {
	locker.LocalLock()
	defer locker.LocalUnlock()
	self.control.Meta["build"] = "true"

	// setup the environment (boot vm)
	// but do not run the dev setup because
	// we dont need any o fthe platform services
	if err := Run("provider_setup", self.control); err != nil {
		return err
	}

	// build code
	if err := Run("code_build", self.control); err != nil {
		return err
	}

	return nil
}
