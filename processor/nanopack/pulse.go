package nanopack

import (
	"github.com/nanobox-io/nanobox/processor"
)

type updatePulse struct {
	config processor.ProcessConfig
}

func init() {
	processor.Register("update_pulse", updatePulseFunc)
}

func updatePulseFunc(config processor.ProcessConfig) (processor.Processor, error) {
	// confirm the provider is an accessable one that we support.

	return updatePulse{config}, nil
}

func (self updatePulse) Results() processor.ProcessConfig {
	return self.config
}

func (self updatePulse) Process() error {
	// TODO: setup the nanoagent services
	return nil
}