package commands

import (
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/processors"
	"github.com/nanobox-io/nanobox/util/display"
)

var (

	// LogoutCmd ...
	LogoutCmd = &cobra.Command{
		Use:   "logout",
		Short: "Removes your nanobox.io api token from your local nanobox client.",
		Long:  ``,
		Run:   logoutFn,
	}
	
	// loginCmdFlags ...
	logoutCmdFlags = struct {
		endpoint string
	}{}
)

func init() {
	LoginCmd.Flags().StringVarP(&logoutCmdFlags.endpoint, "endpoint", "e", "", "endpoint")
}

// logoutFn ...
func logoutFn(ccmd *cobra.Command, args []string) {
	// set default endpoint to nanobox
	if logoutCmdFlags.endpoint == "" {
		logoutCmdFlags.endpoint = "nanobox"
	}
	
	display.CommandErr(processors.Logout(logoutCmdFlags.endpoint))
}
