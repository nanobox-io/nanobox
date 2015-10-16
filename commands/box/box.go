// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

//
package box

import (
	"github.com/nanobox-io/nanobox-cli/config"
	"github.com/nanobox-io/nanobox-cli/util"
	"github.com/nanobox-io/nanobox-cli/util/vagrant"
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
