// package box ...
package box

import (
	"fmt"
	"github.com/nanobox-io/nanobox/config"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/vagrant"
	"github.com/spf13/cobra"
	"os"
)

var (
	BoxCmd = &cobra.Command{
		Use:   "box",
		Short: "Subcommands for managing the nanobox/boot2docker.box",
		Long:  ``,

		PersistentPreRun: runnable,
	}

	Vagrant = vagrant.Default
	Config  = config.Default
	Util    = util.Default
)

//
func init() {
	BoxCmd.AddCommand(installCmd)
	BoxCmd.AddCommand(updateCmd)
}

// runnable ensures all dependencies are satisfied before running box commands
func runnable(ccmd *cobra.Command, args []string) {

	// ensure vagrant exists
	if exists := Vagrant.Exists(); !exists {
		fmt.Println("Missing dependency 'Vagrant'. Please download and install it to continue (https://www.vagrantup.com/).")
		os.Exit(1)
	}

	// ensure virtualbox exists
	if exists := util.VboxExists(); !exists {
		fmt.Println("Missing dependency 'Virtualbox'. Please download and install it to continue (https://www.virtualbox.org/wiki/Downloads).")
		os.Exit(1)
	}
}
