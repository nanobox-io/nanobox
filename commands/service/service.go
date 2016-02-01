// package service ...
package service

import (
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/config"
)

//
var (

	//
	ServiceCmd = &cobra.Command{
		Use:   "service",
		Short: "Subcommands to aid in developing a custom services",
		Long:  ``,
	}

	//
	Config = config.Default

	//
	fFile string // destination file when fetching an engine
)

//
func init() {
	ServiceCmd.AddCommand(fetchCmd)
	ServiceCmd.AddCommand(newCmd)
	ServiceCmd.AddCommand(publishCmd)
}
