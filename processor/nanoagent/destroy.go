package provider

import (
	"github.com/nanobox-io/nanobox/processor"
)

type nanoagentDestroy struct {
	config processor.ProcessConfig
}


func nanoagentDestroyFunc(config processor.ProcessConfig) (Sequence, error) {
	// confirm the provider is an accessable one that we support.

	return nanoagentDestroy{config}, nil
}


func (self nanoagentDestroy) Results() processor.ProcessConfig {
	return self.config
}

func (self nanoagentDestroy) Process() error {
	// TODO: setup the nanoagent services
}