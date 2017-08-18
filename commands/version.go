package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/models"
)

var (

	// VersionCmd prints the nanobox version.
	VersionCmd = &cobra.Command{
		Use:              "version",
		Short:            "Show the current Nanobox version.",
		Long:             ``,
		PersistentPreRun: func(ccmd *cobra.Command, args []string) {},
		Run:              versionFn,
	}
)

// versionFn does the actual printing
func versionFn(ccmd *cobra.Command, args []string) {
	fmt.Println(models.VersionString())
}
