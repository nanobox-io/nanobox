//
package box

import "github.com/spf13/cobra"

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Downloads and adds nanobox/boot2docker box",
	Long:  ``,

	Run: Install,
}

// Install
func Install(ccmd *cobra.Command, args []string) {
	Vagrant.Install()
}
