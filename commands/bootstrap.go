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

	"github.com/nanobox-io/nanobox-cli/config"
	"github.com/nanobox-io/nanobox-cli/util"
	"github.com/nanobox-io/nanobox-golang-stylish"
)

//
var bootstrapCmd = &cobra.Command{
	Hidden: true,

	Use:   "bootstrap",
	Short: "Runs an engine's bootstrap script - downloads code & launches VM",
	Long: `
Description:
  Runs an engine's bootstrap script - downloads code & launches VM`,

	PreRun: bootVM,
	Run:    nanoBootstrap,
}

//
func nanoBootstrap(ccmd *cobra.Command, args []string) {

	// PreRun: bootVM

	fmt.Printf(stylish.Bullet("Bootstrapping code..."))

	//
	bootstrap := util.Sync{
		Model:   "bootstrap",
		Path:    fmt.Sprintf("%s/bootstrap", config.ServerURL),
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
