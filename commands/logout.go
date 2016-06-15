package commands

import (
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/util/print"
)

var (

	// LogoutCmd ...
	LogoutCmd = &cobra.Command{
		Use:   "logout",
		Short: "Removes your nanobox.io api token from your local nanobox client.",
		Long:  ``,
		Run:   logoutFn,
	}
)

// logoutFn ...
func logoutFn(ccmd *cobra.Command, args []string) {
	print.OutputCmdErr(processor.Run("logout", processor.DefaultConfig))
}
