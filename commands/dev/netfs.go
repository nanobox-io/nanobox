package dev

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/util/print"
)

var (

	// NetfsCmd ...
	NetfsCmd = &cobra.Command{
		Hidden: true,
		Use:    "netfs",
		Short:  "Add or remove netfs directories.",
		Long:   ``,
	}

	// NetfsAddCmd ...
	NetfsAddCmd = &cobra.Command{
		Hidden: true,
		Use:    "add",
		Short:  "Add a netfs export.",
		Long:   ``,
		Run:    netfsAddFn,
	}

	// NetfsRmCmd ...
	NetfsRmCmd = &cobra.Command{
		Hidden: true,
		Use:    "rm",
		Short:  "Remove a netfs export.",
		Long:   ``,
		Run:    netfsRmFn,
	}
)

//
func init() {
	NetfsCmd.AddCommand(NetfsAddCmd)
	NetfsCmd.AddCommand(NetfsRmCmd)
}

// netfsAddFn will run the netfs processor for adding a netfs export
func netfsAddFn(ccmd *cobra.Command, args []string) {

	// validate we have args required to set the meta we'll need; if we don't have
	// the required args this will os.Exit(1) with an error message
	if len(args) != 1 {
		fmt.Printf(`
Wrong number of arguments (expecting 1 got %v). Run the command again with the
path of the exports entry you would like to add:

ex: nanobox dev netfs add <path>

`, len(args))
		return
	}

	processor.DefaultConfig.Meta["path"] = args[0]
	print.OutputCommandErr(processor.Run("dev_netfs_add", processor.DefaultConfig))
}

// netfsRmFn will run the netfs processor for removing a netfs export
func netfsRmFn(ccmd *cobra.Command, args []string) {

	// validate we have args required to set the meta we'll need; if we don't have
	// the required args this will os.Exit(1) with an error message
	if len(args) != 1 {
		fmt.Printf(`
Wrong number of arguments (expecting 1 got %v). Run the command again with the
path of the exports entry you would like to remove:

ex: nanobox dev netfs rm <path>

`, len(args))
		return
	}

	// set the meta arguments to be used in the processor and run the processor
	processor.DefaultConfig.Meta["path"] = args[0]
	print.OutputCommandErr(processor.Run("dev_netfs_remove", processor.DefaultConfig))
}
