package nanopack

import (
	"github.com/nanobox-io/nanobox/processor"
)

type nanopackDestroy struct {
	config processor.ProcessConfig
}

func init() {
	processor.Register("nanopack_destroy", nanopackDestroyFunc)
}

func nanopackDestroyFunc(config processor.ProcessConfig) (processor.Processor, error) {
	// confirm the provider is an accessable one that we support.
	return nanopackDestroy{config}, nil
}

func (self nanopackDestroy) Results() processor.ProcessConfig {
	return self.config
}

func (self nanopackDestroy) Process() error {
	// TODO: setup the nanoagent services
	return nil
}
