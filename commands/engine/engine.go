// package engine ...
package engine

import (
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/config"
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
	Config = config.Default

	//
	fFile string // destination file when fetching an engine
)

//
func init() {
	EngineCmd.AddCommand(fetchCmd)
	EngineCmd.AddCommand(newCmd)
	EngineCmd.AddCommand(publishCmd)
}
