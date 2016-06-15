package dev

import (
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/processor"
	cmdutil "github.com/nanobox-io/nanobox/validate/commands"
)

var (

	// InfoCmd ...
	InfoCmd = &cobra.Command{
		Use:    "info",
		Short:  "Displays information about the running app and its components.",
		Long:   ``,
		PreRun: cmdutil.Validate("provider"),
		Run:    infoFn,
	}
)

// infoFn will run the DNS processor for adding DNS entires to the "hosts"
// file
func infoFn(ccmd *cobra.Command, args []string) {
	processor.DefaultConfig.Meta["name"] = args[0]

	//
	if err := processor.Run("dev_dns_add", processor.DefaultConfig); err != nil {

	}
}
