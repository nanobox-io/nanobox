package provider

import (
	"github.com/nanobox-io/nanobox/processor"
)

type updatePulse struct {
	config processor.ProcessConfig
}


func updatePulseFunc(config processor.ProcessConfig) (Sequence, error) {
	// confirm the provider is an accessable one that we support.

	return updatePulse{config}, nil
}


func (self updatePulse) Results() processor.ProcessConfig {
	return self.config
}

func (self updatePulse) Process() error {
	// TODO: setup the nanoagent services
}