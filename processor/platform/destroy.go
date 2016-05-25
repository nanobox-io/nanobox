package platform

import (
	"github.com/nanobox-io/nanobox/processor"
)

type platformDestroy struct {
	config processor.ProcessConfig
}

func init() {
	processor.Register("nanopack_destroy", platformDestroyFunc)
}

func platformDestroyFunc(config processor.ProcessConfig) (processor.Processor, error) {
	// confirm the provider is an accessable one that we support.
	return platformDestroy{config}, nil
}

func (self platformDestroy) Results() processor.ProcessConfig {
	return self.config
}

func (self platformDestroy) Process() error {
	// TODO: setup the nanoagent services
	return nil
}
