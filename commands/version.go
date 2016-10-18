package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	// version info
	tag    string
	commit string

	// versionCmd prints version info and exits
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print version info and exit",
		Run:   versionFn,
	}
)

// versionFn prints version info and exits
func versionFn(ccmd *cobra.Command, args []string) {
	fmt.Printf("nanobox %s (%s)\n", tag, commit)
}
