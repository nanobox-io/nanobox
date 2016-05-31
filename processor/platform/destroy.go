package platform

import (
	"github.com/nanobox-io/nanobox/processor"
)

type platformDestroy struct {
	control processor.ProcessControl
}

func init() {
	processor.Register("nanopack_destroy", platformDestroyFunc)
}

func platformDestroyFunc(control processor.ProcessControl) (processor.Processor, error) {
	// confirm the provider is an accessable one that we support.
	return platformDestroy{control}, nil
}

func (self platformDestroy) Results() processor.ProcessControl {
	return self.control
}

func (self platformDestroy) Process() error {
	// TODO: destroy nanoagent services
	return nil
}
