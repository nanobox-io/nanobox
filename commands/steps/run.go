package steps

import (
	"github.com/nanobox-io/nanobox/commands/registry"
	// "github.com/nanobox-io/nanobox/util/display"
	"github.com/spf13/cobra"
)

func Run(stepNames ...string) func(ccmd *cobra.Command, args []string) {
	//
	return func(ccmd *cobra.Command, args []string) {

		if registry.GetBool("internal") {
			return
		}
		// list that needs to be run
		steps := []step{}
		//
		for _, stepName := range stepNames {
			step, ok := stepList[stepName]
			if ok && !step.complete() {
				steps = append(steps, step)
			}
		}

		// run the missing steps
		for _, step := range steps {
			step.cmd(ccmd, args)
		}
	}
}
