// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

package commands

//
import (
	"github.com/spf13/cobra"

	// "github.com/pagodabox/nanobox-cli/config"
	// "github.com/pagodabox/nanobox-cli/util"
	// "github.com/pagodabox/nanobox-golang-stylish"
)

//
var boxInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "",
	Long:  ``,

	Run: boxInstall,
}

// boxInstall
func boxInstall(ccmd *cobra.Command, args []string) {

}
