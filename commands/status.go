package commands

import (
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/processors"
	"github.com/nanobox-io/nanobox/util/display"
)

var (

	// StatusCmd ...
	StatusCmd = &cobra.Command{
		Use:   "status",
		Short: "Displays the status of your Nanobox VM & running platforms.",
		Long:  ``,
		Run:   statusFn,
	}
)

func statusFn(ccmd *cobra.Command, args []string) {
	display.CommandErr(processors.Status())
}
