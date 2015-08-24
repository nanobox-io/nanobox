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
	"net/url"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/pagodabox/nanobox-cli/config"
	"github.com/pagodabox/nanobox-cli/utils"
	"github.com/pagodabox/nanobox-golang-stylish"
)

//
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Issues a deploy to the nanobox VM",
	Long: `
Description:
  Issues a deploy to the nanobox VM`,

	Run: nanoDeploy,
}

//
func init() {
	deployCmd.Flags().BoolVarP(&fReset, "reset", "r", false, "Clears cached libraries the project might use")
	deployCmd.Flags().BoolVarP(&fSandbox, "sandbox", "s", false, "Creates your app environment w/o webs or workers")
	deployCmd.Flags().BoolVarP(&fVerbose, "verbose", "v", false, "Increases the level of log output from 'info' to 'debug'")
}

// nanoDeploy
func nanoDeploy(ccmd *cobra.Command, args []string) {
	fmt.Printf(stylish.Bullet("Deploying codebase..."))

	v := url.Values{}

	v.Add("reset", strconv.FormatBool(fReset))
	v.Add("sandbox", strconv.FormatBool(fSandbox))

	//
	deploy := utils.Sync{
		Model:   "deploy",
		Path:    fmt.Sprintf("http://%v:1757/deploys?%v", config.Nanofile.IP, v.Encode()),
		Verbose: fVerbose,
	}

	//
	deploy.Run(args)

	//
	switch deploy.Status {

	// complete
	case "complete":

		//
		if fSandbox {
			fmt.Printf(stylish.Bullet("Sandbox deploy complete..."))
			break
		}

		fmt.Printf(stylish.Bullet(fmt.Sprintf("Deploy complete... Navigate to %v.nano.dev to view your app.", config.App)))

		// if the deploy fails the server should handle the message. If not, this can
		// be re-enabled
	case "errored":
		// fmt.Printf(stylish.Error("Deploy failed", "Your deploy failed to well... deploy"))
	}
}
