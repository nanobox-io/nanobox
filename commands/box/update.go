//
package box

import "github.com/spf13/cobra"

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Updates nanobox/boot2docker box",
	Long:  ``,

	Run: Update,
}

// Update
func Update(ccmd *cobra.Command, args []string) {
	Vagrant.Update()
}
