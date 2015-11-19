//
package commands

import (
	"fmt"
	"github.com/nanobox-io/nanobox-golang-stylish"
	"github.com/spf13/cobra"
)

//
var sshCmd = &cobra.Command{
	Hidden: true,

	Use:   "ssh",
	Short: "SSH into the nanobox",
	Long:  ``,

	PreRun: boot,
	Run:    ssh,
}

// ssh
func ssh(ccmd *cobra.Command, args []string) {

	// PreRun: boot

	fmt.Printf(stylish.Bullet("SSHing into nanobox..."))
	if err := Vagrant.SSH(); err != nil {
		Config.Fatal("[commands/ssh] vagrant.SSH() failed", err.Error())
	}
}
