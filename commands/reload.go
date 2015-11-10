//
package commands

import (
	"fmt"
	"github.com/nanobox-io/nanobox-golang-stylish"
	"github.com/spf13/cobra"
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

	fmt.Printf(stylish.Bullet("Reloading nanobox..."))
	if err := Vagrant.Reload(); err != nil {
		Config.Fatal("[commands/reload] vagrant.Reload() failed - ", err.Error())
	}
}
