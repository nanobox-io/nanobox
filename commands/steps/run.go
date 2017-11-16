package steps

import (
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/commands/registry"
)

func Run(stepNames ...string) func(ccmd *cobra.Command, args []string) {
	return func(ccmd *cobra.Command, args []string) {

		if registry.GetBool("internal") {
			return
		}

		for _, stepName := range stepNames {
			step, ok := stepList[stepName]
			if ok && !step.complete() {
				step.cmd(ccmd, args)
			}
		}
	}
}
