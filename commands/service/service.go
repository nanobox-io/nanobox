// package service ...
package service

import (
	"github.com/spf13/cobra"
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
	fFile string // destination file when fetching an engine
)

//
func init() {
	ServiceCmd.AddCommand(fetchCmd)
	ServiceCmd.AddCommand(newCmd)
	ServiceCmd.AddCommand(publishCmd)
}
