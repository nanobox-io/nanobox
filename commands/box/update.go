//
package box

import (
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/util/vagrant"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Updates nanobox/boot2docker box",
	Long:  ``,

	Run: Update,
}

// Update
func Update(ccmd *cobra.Command, args []string) {
	vagrant.Update()
}
