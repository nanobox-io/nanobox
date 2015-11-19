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
	Hidden: true,

	Use:   "reload",
	Short: "Reloads the nanobox",
	Long:  ``,

	PreRun: initialize,
	Run:    reload,
}

// reload runs 'vagrant reload --provision'
func reload(ccmd *cobra.Command, args []string) {

	// PreRun: initialize

	//
	fmt.Printf(stylish.Bullet("Reloading nanobox..."))
	fmt.Printf(stylish.Bullet("Nanobox may require admin privileges to modify your /etc/hosts file and /etc/exports."))
	if err := vagrant.Reload(); err != nil {
		vagrant.Fatal("[commands/reload] vagrant.Reload() failed", err.Error())
	}

	// indeicate that the VM has recently been reloaded
	config.VMfile.ReloadedIs(true)
}
