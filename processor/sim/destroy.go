
package sim

import (
	"fmt"

	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/util/config"
)

// processSimDestroy ...
type processSimDestroy struct {
	control processor.ProcessControl
}

//
func init() {
	processor.Register("sim_destroy", simDestroyFn)
}

//
func simDestroyFn(control processor.ProcessControl) (processor.Processor, error) {
	simDestroy := &processSimDestroy{control}
	return simDestroy, simDestroy.validateMeta()
}

func (simDestroy *processSimDestroy) validateMeta() error {
	if simDestroy.control.Meta["app_name"] == "" {
		simDestroy.control.Meta["app_name"] = fmt.Sprintf("%s_sim", config.AppID())
	}

	return nil
}

//
func (simDestroy processSimDestroy) Results() processor.ProcessControl {
	return simDestroy.control
}

//
func (simDestroy *processSimDestroy) Process() error {
	simDestroy.control.Env = "sim"

	return processor.Run("env_destroy", simDestroy.control)
}
