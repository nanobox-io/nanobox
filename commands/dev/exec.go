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
		PreRun: validate.Requires("provider", "provider_up", "built", "dev_deployed"),
		Run:    execFn,
	}

	// execCmdFlags ...
	execCmdFlags = struct {
		container string
	}{}
)

func init() {
	ExecCmd.Flags().StringVarP(&execCmdFlags.container, "container", "c", "", "specify the container in which to exec the command")
}

// execFn ...
func execFn(ccmd *cobra.Command, args []string) {

	// validate we have args required to set the meta we'll need; if we don't have
	// the required args this will return with instructions
	if len(args) < 1 {
		fmt.Printf(`
Wrong number of arguments (expecting 1 or more got %v). Run the command again
with the command you would like to exec:

ex: nanobox dev exec <command>

	`, len(args))

		return
	}

	// set the meta arguments to be used in the processor and run the processor
	processor.DefaultControl.Meta["command"] = strings.Join(args[0:], " ")
	processor.DefaultControl.Meta["container"] = execCmdFlags.container
	print.OutputCommandErr(processor.Run("dev_console", processor.DefaultControl))
}
