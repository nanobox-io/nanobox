package commands

import (
  "github.com/spf13/cobra"

  "github.com/nanobox-io/nanobox/util/netfs"
)

var (
  DevNetfsCmd = &cobra.Command{
    Use: "netfs",
    Short: "add or remove netfs directories",
    Long: ``,
    // Hidden: true,
  }

  DevNetfsAddCmd = &cobra.Command{
    Use: "add",
    Short: "add a netfs export",
    Long: ``,
    // Hidden: true,
    Run: devNetfsAddFunc,
  }

  DevNetfsRmCmd = &cobra.Command{
    Use: "rm",
    Short: "remove a netfs export",
    Long: ``,
    // Hidden: true,
    Run: devNetfsRmFunc,
  }

)

func init() {
  DevNetfsCmd.AddCommand(DevNetfsAddCmd)
  DevNetfsCmd.AddCommand(DevNetfsRmCmd)
}

// devNetfsAddFunc will run the netfs function for adding a netfs export
func devNetfsAddFunc(ccmd *cobra.Command, args[]string) {
  // validate that a path was provided
  path := args[0]

  // todo: error if path is nil

  netfs.Add(path)
}

// devNetfsRmFunc will run the netfs function for removing a netfs export
func devNetfsRmFunc(ccmd *cobra.Command, args[]string) {
  // validate that a path was provided
  path := args[0]

  // todo: error if path is nil

  netfs.Remove(path)
}
