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
		notRun := []string{}

		// 
		for _, stepName := range stepNames {
			step, ok := stepList[stepName]
			if ok && !step.complete()	{
				notRun = append(notRun, stepName)
			}
		}

		if len(notRun) == 0 {
			return
		}

		// print the message if 
		printMessage(notRun)

		// run the missing steps
		for _, nr := range notRun {
			step := stepList[nr]

			display.OpenContext("(%s %s)", os.Args[0], nr)
			step.cmd(ccmd, args)
			display.CloseContext()
		}
	}
}


func printMessage(notRun []string) {
	fmt.Println(`
NOTE --------------------------------------------------
Before we can run that command, we need to run these
commands first (generic message):
`)
	for _, nr := range notRun {
		fmt.Printf("$ %s %s\n", os.Args[0], nr)
	}

	fmt.Println("-------------------------------------------------------")
}

