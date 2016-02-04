//
package dev

import (
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/config"
	"github.com/nanobox-io/nanobox/util/vagrant"
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
	if err := vagrant.Suspend(); err != nil {
		vagrant.Fatal("[commands/stop] vagrant.Suspend() failed", err.Error())
	}

	// boot the machine normally (not backgrounded) on next command
	config.VMfile.BackgroundIs(false)
}
