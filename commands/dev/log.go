package dev

import (
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/util/print"
	"github.com/nanobox-io/nanobox/validate"
)

var (

	// LogCmd ...
	LogCmd = &cobra.Command{
		Use:    "log",
		Short:  "Displays logs from the running dev app and its components.",
		Long:   ``,
		PreRun: validate.Requires("provider", "provider_up", "built", "dev_isup"),
		Run:    logFn,
	}
)

// logFn will run the DNS processor for adding DNS entires to the "hosts" file
func logFn(ccmd *cobra.Command, args []string) {
	print.OutputCommandErr(processor.Run("dev_log", processor.DefaultControl))
}
