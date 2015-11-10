// package engine ...
package engine

import (
	"github.com/nanobox-io/nanobox/auth"
	"github.com/nanobox-io/nanobox/config"
	// "github.com/nanobox-io/nanobox/util/s3"
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
	Config = config.Default
	Auth   = auth.Default
	// S3     = s3.Default

	//
	fFile string // destination file when fetching an engine
)

//
func init() {
	EngineCmd.AddCommand(fetchCmd)
	EngineCmd.AddCommand(newCmd)
	EngineCmd.AddCommand(publishCmd)
}
