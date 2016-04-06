package nanoagent

import (
	"github.com/nanobox-io/nanobox/processor"
)

func init() {
	processor.Register("nanoagent_setup", nanoagentSetupFunc)
	processor.Register("nanoagent_destroy", nanoagentDestroyFunc)
	processor.Register("update_pulse", updatePulseFunc)
}
