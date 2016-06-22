package share

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/jcelliott/lumber"
	"github.com/nanobox-io/golang-docker-client"
	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/util/config"
)

// processSimConsole ...
type processSimConsole struct {
	control   processor.ProcessControl
}

// the sim console is the same as the share console
func init() {
	processor.Register("share_console", simConsoleFn)
}

//
func simConsoleFn(control processor.ProcessControl) (processor.Processor, error) {
	simConsole := &processShareSim{control: control}
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

