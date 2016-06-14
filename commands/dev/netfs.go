package dev

import (
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/processor"
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
		Run:    netfsAddFunc,
	}

	// NetfsRmCmd ...
	NetfsRmCmd = &cobra.Command{
		Hidden: true,
		Use:    "rm",
		Short:  "Remove a netfs export.",
		Long:   ``,
		Run:    netfsRmFunc,
	}
)

//
func init() {
	NetfsCmd.AddCommand(NetfsAddCmd)
	NetfsCmd.AddCommand(NetfsRmCmd)
}

// netfsAddFunc will run the netfs processor for adding a netfs export
func netfsAddFunc(ccmd *cobra.Command, args []string) {
	processor.DefaultConfig.Meta["host"] = args[0]
	processor.DefaultConfig.Meta["path"] = args[1]
	processor.Run("dev_netfs_add", processor.DefaultConfig)
}

// netfsRmFunc will run the netfs processor for removing a netfs export
func netfsRmFunc(ccmd *cobra.Command, args []string) {
	processor.DefaultConfig.Meta["host"] = args[0]
	processor.DefaultConfig.Meta["path"] = args[1]
	processor.Run("dev_netfs_remove", processor.DefaultConfig)
}
