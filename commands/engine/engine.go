// package engine ...
package engine

import (
	"github.com/spf13/cobra"
)

//
var (

	//
	EngineCmd = &cobra.Command{
		Use:   "engine",
		Short: "Subcommands to aid in developing a custom engine",
		Long:  ``,
	}

	//
	fFile string // destination file when fetching an engine
)

//
func init() {
	EngineCmd.AddCommand(fetchCmd)
	EngineCmd.AddCommand(newCmd)
	EngineCmd.AddCommand(publishCmd)
}
