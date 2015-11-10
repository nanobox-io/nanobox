//
package main

import (
	"fmt"
	"github.com/nanobox-io/nanobox/commands"
	"github.com/nanobox-io/nanobox/config"
	"os"
	"os/exec"
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

	pass := true

	// ensure vagrant is installed
	if err := exec.Command("vagrant", "-v").Run(); err != nil {
		fmt.Println("Missing dependency 'Vagrant'. Please download and install it to continue (https://www.vagrantup.com/).")
		pass = false
	}

	// ensure virtualbox is installed
	if err := exec.Command("vboxmanage", "-v").Run(); err != nil {
		fmt.Println("Missing dependency 'Virtualbox'. Please download and install it to continue (https://www.virtualbox.org/wiki/Downloads).")
		pass = false
	}

	// if a dependency check fails, exit
	if !pass {
		return
	}

	// check to see if nanobox needs to be updated
	commands.Update()

	//
	commands.NanoboxCmd.Execute()
}
