package provider

import (
	"github.com/nanobox-io/nanobox/processor"
)

type nanoagentSetup struct {
	config processor.ProcessConfig
}


func nanoagentSetupFunc(config processor.ProcessConfig) (Sequence, error) {
	// confirm the provider is an accessable one that we support.

	return nanoagentSetup{config}, nil
}


func (self nanoagentSetup) Results() processor.ProcessConfig {
	return self.config
}

func (self nanoagentSetup) Process() error {
	// TODO: setup the nanoagent services
}