//
package commands

import (
	"fmt"

	"github.com/nanobox-io/nanobox-golang-stylish"
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/config"
	"github.com/nanobox-io/nanobox/util/vagrant"
)

//
var reloadCmd = &cobra.Command{
	Use:   "reload",
	Short: "Reloads the nanobox",
	Long:  ``,

	PreRun: initialize,
	Run:    reload,
}

// reload runs 'vagrant reload --provision'
func reload(ccmd *cobra.Command, args []string) {

	// PreRun: initialize

	// generate a new Vagrantfile on reload; this is done because there may be times
	// when a user needs a new Vagrantfile (ie. adding a custom engine after the VM
	// already exists). The only way to accomplish this otherwise is either destroying
	// the VM entirely or running the hidden "nanobox init -f"
	Vagrant.Init()

	//
	fmt.Printf(stylish.Bullet("Reloading nanobox..."))
	fmt.Printf(stylish.Bullet("Nanobox may require admin privileges to modify your /etc/hosts and /etc/exports files."))
	if err := vagrant.Reload(); err != nil {
		vagrant.Fatal("[commands/reload] vagrant.Reload() failed", err.Error())
	}

	// indicate that the VM has recently been reloaded; we do this under the assumption
	// that the reload command was called because something went wrong with the VM
	// and a user is trying to "reset" things.
	config.VMfile.ReloadedIs(true)
}
