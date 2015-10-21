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
	"github.com/nanobox-io/nanobox-golang-stylish"
	"github.com/nanobox-io/nanobox/config"
	"github.com/spf13/cobra"
	"net/url"
	"strconv"
)

var (

	//
	deployCmd = &cobra.Command{
		Hidden: true,

		Use:   "deploy",
		Short: "Deploys code to the nanobox",
		Long:  ``,

		PreRun:  boot,
		Run:     deploy,
		PostRun: halt,
	}

	//
	install bool // tells nanobox server to install services
)

//
func init() {
	deployCmd.Flags().BoolVarP(&install, "run", "", false, "Creates your app environment w/o webs or workers")
}

// deploy
func deploy(ccmd *cobra.Command, args []string) {

	// PreRun: boot

	fmt.Printf(stylish.Bullet("Deploying codebase..."))

	// stream deploy output
	go Mist.Stream([]string{"log", "deploy"}, Mist.PrintLogStream)

	// listen for status updates
	errch := make(chan error)
	go func() {
		errch <- Mist.Listen([]string{"job", "deploy"}, Mist.DeployUpdates)
	}()

	v := url.Values{}
	v.Add("reset", strconv.FormatBool(config.Force))
	v.Add("run", strconv.FormatBool(install))

	// run a deploy
	if err := Server.Deploy(v.Encode()); err != nil {
		Config.Fatal("[commands/deploy] server.Deploy() failed - ", err.Error())
	}

	// wait for a status update (blocking)
	err := <-errch

	//
	if err != nil {
		fmt.Printf(err.Error())
		return
	}

	// PostRun: halt
}
