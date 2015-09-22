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
	"github.com/pagodabox/nanobox-cli/util"
	"github.com/pagodabox/nanobox-golang-stylish"
)

//
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Issues a deploy to the nanobox VM",
	Long: `
Description:
  Issues a deploy to the nanobox VM

  -f, --force[=false]: Clears cached libraries the project might use`,

	PreRun: VMIsRunning,
	Run:    nanoDeploy,
}

//
func init() {
	deployCmd.Flags().BoolVarP(&fSandbox, "sandbox", "s", false, "Creates your app environment w/o webs or workers")
}

// nanoDeploy
func nanoDeploy(ccmd *cobra.Command, args []string) {
	fmt.Printf(stylish.Bullet("Deploying codebase..."))

	v := url.Values{}

	v.Add("reset", strconv.FormatBool(fForce))
	v.Add("sandbox", strconv.FormatBool(fSandbox))

	//
	deploy := util.Sync{
		Model:   "deploy",
		Path:    fmt.Sprintf("http://%s/deploys?%v", config.ServerURI, v.Encode()),
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

		fmt.Printf(stylish.Bullet("Deploy complete... Navigate to %v.nano.dev to view your app.", config.App))

		// if the deploy fails the server should handle the message. If not, this can
		// be re-enabled
	case "errored":
		// fmt.Printf(stylish.Error("Deploy failed", "Your deploy failed to well... deploy"))
	}
}
