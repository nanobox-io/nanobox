package dev

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/util/print"
	"github.com/nanobox-io/nanobox/validate"
)

var (

	// ExecCmd ...
	ExecCmd = &cobra.Command{
		Use:    "exec",
		Short:  "Executes a command inside your local dev app.",
		Long:   ``,
		PreRun: validate.Requires("provider"),
		Run:    execFn,
	}
)

// execFn ...
func execFn(ccmd *cobra.Command, args []string) {

	// validate we have args required to set the meta we'll need; if we don't have
	// the required args this will os.Exit(1) with instructions
	switch {

	// if no arguments are passed or too many arguments are passed let the user
	// know they are using the command wrong
	case len(args) == 0:
		fmt.Printf(`
Wrong number of arguments (expecting 1 or more got %v). Run the command again
with the command you would like to exec OR the name of the container you would
like the command to exec inside of and the command to exec:

ex: nanobox dev exec <command>
ex: nanobox dev exec <container> <command>

	`, len(args))
		return

	// if one argument is passed we'll assume that it's the command trying to exec
	// in the "default" container
	case len(args) == 1:
		processor.DefaultConfig.Meta["command"] = args[0]

	// if 2 arguemnts are passed we'll assume the first is the container and the
	// second one is the command
	case len(args) == 2:
		processor.DefaultConfig.Meta["container"] = args[0]
		processor.DefaultConfig.Meta["command"] = strings.Join(args[1:], " ")
	}

	// set the meta arguments to be used in the processor and run the processor
	print.OutputCommandErr(processor.Run("dev_console", processor.DefaultConfig))
}
