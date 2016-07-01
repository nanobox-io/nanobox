package sim

import (
	"github.com/nanobox-io/nanobox/processor"
)

// processSimLog ...
type processSimLog struct {
	control processor.ProcessControl
}

//
func init() {
	processor.Register("sim_log", simLogFn)
}

//
func simLogFn(control processor.ProcessControl) (processor.Processor, error) {
	return processSimLog{control}, nil
}

//
func (simLog processSimLog) Results() processor.ProcessControl {
	return simLog.control
}

//
func (simLog processSimLog) Process() error {
	// set the process mode to sim
	simLog.control.Env = "sim"

	// some messaging about the logging??

	return processor.Run("mist_log", simLog.control)
}
