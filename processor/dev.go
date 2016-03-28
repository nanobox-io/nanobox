package processor

import (
)

type dev struct {
	config ProcessConfig
}

func init() {
	Register("dev", devFunc)
}

func devFunc(config ProcessConfig) (Sequence, error) {
	config.Meta["dev-config"]
	// do some config validation
	// check on the meta for the flags and make sure they work

	return dev{config}, nil
}

func (self dev) Results() ProcessConfig {
	return self.config
}

func (self dev) Process() error {
	// setup the environment (boot vm)

	// setup remote api (make sure nanoagent is working)

	// setup mist

	// setup portal

	// setup Logvac

	// build code

	// start services

	// start code

	// update portal

	return nil
}