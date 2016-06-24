package sim

import (
	"github.com/nanobox-io/nanobox/processor"
)

// processSimConsole ...
type processSimConsole struct {
	control   processor.ProcessControl
}

// the sim console is the same as the share console
func init() {
	processor.Register("sim_console", simConsoleFn)
}

//
func simConsoleFn(control processor.ProcessControl) (processor.Processor, error) {
	simConsole := &processSimConsole{control: control}
	return simConsole, nil
}

func (simConsole processSimConsole) Results() processor.ProcessControl {
	return simConsole.control
}

// this process is just a shortcut so we can do any other special
// stuff. Which currently there is nothing other tuen running
// the share console.
func (simConsole processSimConsole) Process() error {
	return processor.Run("share_console", simConsole.control)
}

