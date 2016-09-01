package steps

import (
	"os"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/nanobox-io/nanobox/util/display"

) 

func Run(stepNames ...string) func(ccmd *cobra.Command, args []string) {
	//
	return func(ccmd *cobra.Command, args []string) {

		// list that needs to be run
		prereqs := []string{}

		// 
		for _, stepName := range stepNames {
			step, ok := stepList[stepName]
			if ok && !step.complete()	{
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

			display.OpenContext("(%s %s)", os.Args[0], nr)
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
	fmt.Println("You're welcome ;)")
	fmt.Println("------------------------------------")
	fmt.Println()
}
