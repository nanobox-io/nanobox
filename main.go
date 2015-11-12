//
package main

import (
	"fmt"
	"github.com/nanobox-io/nanobox/commands"
	"github.com/nanobox-io/nanobox/config"
	"os"
	"runtime/debug"
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
			config.Log.Fatal(fmt.Sprintf("Cause of failure: %v", r))
			config.Log.Fatal(fmt.Sprintf("Error output:\n%v\n", string(stack)))
			config.Log.Close()
			os.Exit(1)
		}
	}()

	// check to see if nanobox needs to be updated
	commands.Update()

	//
	commands.NanoboxCmd.Execute()
}
