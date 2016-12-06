package env

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/processors/env/share"
	"github.com/nanobox-io/nanobox/util/display"
)

var (

	// ShareCmd ...
	ShareCmd = &cobra.Command{
		Hidden: true,
		Use:    "share",
		Short:  "Add or remove share directories.",
		Long:   ``,
	}

	// ShareAddCmd ...
	ShareAddCmd = &cobra.Command{
		Hidden: true,
		Use:    "add",
		Short:  "Add a share export.",
		Long:   ``,
		Run:    shareAddFn,
	}

	// ShareRmCmd ...
	ShareRmCmd = &cobra.Command{
		Hidden: true,
		Use:    "rm",
		Short:  "Remove a share export.",
		Long:   ``,
		Run:    shareRmFn,
	}
)

//
func init() {
	ShareCmd.AddCommand(ShareAddCmd)
	ShareCmd.AddCommand(ShareRmCmd)
}

// shareAddFn will run the share processor for adding a share export
func shareAddFn(ccmd *cobra.Command, args []string) {

	// validate we have args required to set the meta we'll need; if we don't have
	// the required args this will return with instructions
	if len(args) != 1 {
		fmt.Printf(`
Wrong number of arguments (expecting 1 got %v). Run the command again with the
path of the exports entry you would like to add:

ex: nanobox env share add <path>

`, len(args))

		return
	}

	display.CommandErr(share.Add(args[0]))
}

// shareRmFn will run the share processor for removing a share export
func shareRmFn(ccmd *cobra.Command, args []string) {

	// validate we have args required to set the meta we'll need; if we don't have
	// the required args this will return with instructions
	if len(args) != 1 {
		fmt.Printf(`
Wrong number of arguments (expecting 1 got %v). Run the command again with the
path of the exports entry you would like to remove:

ex: nanobox env share rm <path>

`, len(args))

		return
	}

	display.CommandErr(share.Remove(args[0]))
}
