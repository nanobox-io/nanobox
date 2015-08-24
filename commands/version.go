// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

package commands

//
import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/pagodabox/nanobox-cli/config"
)

//
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Run anytime to see the current version of the CLI.",
	Long: `
Description:
  Run anytime to see the current version of the CLI.`,

	Run: nanoVersion,
}

// nanoVersion runs 'vagrant suspend'
func nanoVersion(ccmd *cobra.Command, args []string) {
	fmt.Printf("Version: %v\n", config.Version.String())
}
