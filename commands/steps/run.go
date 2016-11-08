package steps

import (
	"fmt"

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
		prereqs := []string{}
		steps := []step{}
		//
		for _, stepName := range stepNames {
			step, ok := stepList[stepName]
			if ok && !step.complete() {
				steps = append(steps, step)
				if !step.private {
					prereqs = append(prereqs, stepName)
				}
			}
		}

		if len(steps) == 0 {
			return
		}

		// print the message if
		printMessage(prereqs)

		// run the missing steps
		for _, step := range steps {
			step.cmd(ccmd, args)
		}
	}
}

func printMessage(prereqs []string) {
	if len(prereqs) == 0 {
		return
	}
	fmt.Println()
	fmt.Println("------------------------------------")
	fmt.Println("Running the following prerequisites:")
	fmt.Println()

	for _, dep := range prereqs {
		fmt.Printf("$ nanobox %s\n", dep)
	}

	fmt.Println()
	fmt.Println("------------------------------------")
	fmt.Println()
}
