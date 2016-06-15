package platform

import "github.com/nanobox-io/nanobox/processor"

// processPlatformDestroy ...
type processPlatformDestroy struct {
	control processor.ProcessControl
}

//
func init() {
	processor.Register("nanopack_destroy", platformDestroyFn)
}

//
func platformDestroyFn(control processor.ProcessControl) (processor.Processor, error) {
	return processPlatformDestroy{control}, nil
}

//
func (platformDestroy processPlatformDestroy) Results() processor.ProcessControl {
	return platformDestroy.control
}

// TODO: destroy nanoagent services
func (platformDestroy processPlatformDestroy) Process() error {
	return nil
}
