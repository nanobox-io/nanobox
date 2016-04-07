//
package main

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/jcelliott/lumber"

	"github.com/nanobox-io/nanobox/commands"
	
)

// main
func main() {

	// global panic handler; this is done to avoid showing any panic output if
	// something happens to fail. The output is logged and "pretty" message is
	// shown
	defer func() {
		if r := recover(); r != nil {
			// put r into your log ( it contains the panic message)
			// Then log debug.Stack (from the runtime/debug package)

			stack := debug.Stack()

			fmt.Println("Nanobox encountered an unexpected error. Please see ~/.nanobox/nanobox.log and submit the issue to us.")
			lumber.Fatal(fmt.Sprintf("Cause of failure: %v", r))
			lumber.Fatal(fmt.Sprintf("Error output:\n%v\n", string(stack)))
			lumber.Close()
			os.Exit(1)
		}
	}()

	// commented because we arent doing updates right now.
	// // check to see if nanobox needs to be updated
	// if err := commands.Update(); err != nil {
	// 	fmt.Println("Nanobox was unable to update because of the following error:\n", err.Error())
	// }

	//
	commands.NanoboxCmd.Execute()
}
