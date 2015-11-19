//
package commands

import (
	"github.com/nanobox-io/nanobox/config"
	"github.com/spf13/cobra"
)

//
var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Suspends the nanobox",
	Long:  ``,

	PreRun: initialize,
	Run:    stop,
}

// stop runs 'vagrant suspend'
func stop(ccmd *cobra.Command, args []string) {

	// PreRun: initialize

	//
	if err := Vagrant.Suspend(); err != nil {
		Config.Fatal("[commands/stop] vagrant.Suspend() failed", err.Error())
	}

	// boot the machine normally (not backgrounded) on next command
	config.VMfile.BackgroundIs(false)
}
