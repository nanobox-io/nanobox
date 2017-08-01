package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/models"
)

var (

	// VersionCmd ...
	VersionCmd = &cobra.Command{
		Use:   "version",
		Short: "Show the current Nanobox version.",
		Long:  ``,
		Run:   versionFn,
	}
)

// versionFn ...
func versionFn(ccmd *cobra.Command, args []string) {
	fmt.Println(models.VersionString())
}
