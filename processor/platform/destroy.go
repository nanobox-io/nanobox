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
	// confirm the provider is an accessable one that we support.
	return processPlatformDestroy{control}, nil
}

//
func (platformDestroy processPlatformDestroy) Results() processor.ProcessControl {
	return platformDestroy.control
}

//
func (platformDestroy processPlatformDestroy) Process() error {
	// TODO: destroy nanoagent services
	return nil
}
