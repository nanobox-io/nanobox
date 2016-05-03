//
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"

	"github.com/jcelliott/lumber"

	"github.com/nanobox-io/nanobox/commands"
	"github.com/nanobox-io/nanobox/util"
)

// main
func main() {
	// setup a file logger, this will be replaced in verbose mode.
	fileLogger, err := lumber.NewTruncateLogger(filepath.ToSlash(filepath.Join(util.GlobalDir(), "nanobox.log")))
	if err != nil {
		fmt.Println("logging error:", err)
	}

	lumber.SetLogger(fileLogger)
	lumber.Level(lumber.DEBUG)
	// global panic handler; this is done to avoid showing any panic output if
	// something happens to fail. The output is logged and "pretty" message is
	// shown
	defer func() {
		if r := recover(); r != nil {
			// put r into your log ( it contains the panic message)
			// Then log debug.Stack (from the runtime/debug package)

			stack := debug.Stack()

			lumber.Fatal(fmt.Sprintf("Cause of failure: %v", r))
			lumber.Fatal(fmt.Sprintf("Error output:\n%v\n", string(stack)))
			lumber.Close()
			fmt.Println("Nanobox encountered an unexpected error. Please see ~/.nanobox/nanobox.log and submit the issue to us.")
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
