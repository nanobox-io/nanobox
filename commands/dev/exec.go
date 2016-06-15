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
	switch len(args) {
	case 0:
		fmt.Println("I need atleast one arguement to execute")
		return
	case 1:
		processor.DefaultConfig.Meta["command"] = args[0]
	default:
		processor.DefaultConfig.Meta["name"] = args[0]
		processor.DefaultConfig.Meta["command"] = strings.Join(args[1:], " ")
	}

	//
	print.OutputCommandErr(processor.Run("dev_console", processor.DefaultConfig))
}
