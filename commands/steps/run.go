package steps

import (
	"fmt"
	"os"

	"github.com/nanobox-io/nanobox/commands/registry"
	"github.com/nanobox-io/nanobox/util/display"
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

		//
		for _, stepName := range stepNames {
			step, ok := stepList[stepName]
			if ok && !step.complete() {
				prereqs = append(prereqs, stepName)
			}
		}

		if len(prereqs) == 0 {
			return
		}

		// print the message if
		printMessage(prereqs)

		// run the missing steps
		for _, nr := range prereqs {
			step := stepList[nr]

			display.OpenContext("(nanobox %s)", nr)
			step.cmd(ccmd, args)
			display.CloseContext()
		}
	}
}

func printMessage(prereqs []string) {
	fmt.Println()
	fmt.Println("------------------------------------")
	fmt.Println("Running the following prerequisites:")
	fmt.Println()

	for _, dep := range prereqs {
		fmt.Printf("$ %s %s\n", os.Args[0], dep)
	}

	fmt.Println()
	fmt.Println("------------------------------------")
	fmt.Println()
}
