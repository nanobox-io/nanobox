// package box ...
package box

import (
	"github.com/nanobox-io/nanobox/config"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/vagrant"
	"github.com/spf13/cobra"
)

var (
	BoxCmd = &cobra.Command{
		Use:   "box",
		Short: "Subcommands for managing the nanobox/boot2docker.box",
		Long:  ``,
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
