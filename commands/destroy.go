//
package commands

import (
	"fmt"
	"github.com/nanobox-io/nanobox-golang-stylish"
	"github.com/nanobox-io/nanobox/config"
	"github.com/nanobox-io/nanobox/util/vagrant"
	"github.com/spf13/cobra"
	"os"
)

var (

	//
	destroyCmd = &cobra.Command{
		Use:   "destroy",
		Short: "Destroys the nanobox",
		Long:  ``,

		Run: destroy,
	}

	//
	removeEntry bool // does an entry need to be added to the /etc/hosts file
)

//
func init() {
	destroyCmd.Flags().BoolVarP(&removeEntry, "remove-entry", "", false, "")
	destroyCmd.Flags().MarkHidden("remove-entry")
}

// destroy
func destroy(ccmd *cobra.Command, args []string) {

	// if the command is being run with --remove-entry, it means an entry needs
	// to be removed from the hosts file and execution yielded back to the parent
	if removeEntry {
		Hosts.RemoveDomain()
		os.Exit(0) // this exits the sudoed (child) destroy, not the parent proccess
	}

	// destroy the vm; this needs to happen before cleaning up the app to ensure
	// there is a Vagrantfile to run the command with (otherwise it will just get
	// re-created)
	fmt.Printf(stylish.Bullet("Destroying nanobox..."))
	fmt.Printf(stylish.Bullet("Nanobox may require admin privileges to modify your /etc/hosts file and /etc/exports."))
	if err := vagrant.Destroy(); err != nil {

		// dont care if the project no longer exists... thats what we're doing anyway
		if err != err.(*os.PathError) {
			vagrant.Fatal("[commands/destroy] vagrant.Destroy() failed", err.Error())
		}
	}

	// remove app; this needs to happen after the VM is destroyed so that the app
	// isn't just created again upon running the vagrant command
	fmt.Printf(stylish.Bullet("Deleting nanobox files (%s)", config.AppDir))
	if err := os.RemoveAll(config.AppDir); err != nil {
		Config.Fatal("[commands/destroy] os.RemoveAll() failed", err.Error())
	}

	// attempt to remove the entry regardless of whether its there or not
	sudo("destroy --remove-entry", fmt.Sprintf("Removing %s domain from /etc/hosts", config.Nanofile.Domain))
}
