// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/pagodabox/nanobox-cli/config"
	"github.com/pagodabox/nanobox-cli/utils"
	"github.com/pagodabox/nanobox-golang-stylish"
)

//
var bootstrapCmd = &cobra.Command{
	Use:   "bootstrap",
	Short: "Runs an engine's bootstrap script - downloads code & launches VM",
	Long: `
Description:
  Runs an engine's bootstrap script - downloads code & launches VM`,

	Run: nanoBootstrap,
}

//
func init() {
	bootstrapCmd.Flags().BoolVarP(&fVerbose, "verbose", "v", false, "Increases the level of log output from 'info' to 'debug'")
}

//
func nanoBootstrap(ccmd *cobra.Command, args []string) {
	fmt.Printf(stylish.Bullet("Bootstrapping code..."))

	//
	bootstrap := utils.Sync{
		Model:   "bootstrap",
		Path:    fmt.Sprintf("http://%v:1757/bootstrap", config.Nanofile.IP),
		Verbose: fVerbose,
	}

	//
	bootstrap.Run(args)

	//
	switch bootstrap.Status {

	// complete
	case "complete":
		fmt.Printf(stylish.Bullet("Bootstrap complete"))

	// if the bootstrap fails the server should handle the message. If not, this can
	// be re-enabled
	case "errored":
		// fmt.Printf(stylish.Error("Bootstrap failed", "Your app failed to bootstrap"))
	}
}
