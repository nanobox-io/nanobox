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
	"os"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox-cli/config"
	"github.com/nanobox-io/nanobox-cli/util"
	// "github.com/nanobox-io/nanobox-golang-stylish"
)

//
var deployCmd = &cobra.Command{
	Hidden: true,

	Use:   "deploy",
	Short: "Deploys code to the nanobox",
	Long:  ``,

	PreRun: bootVM,
	Run:    nanoDeploy,
}

//
func init() {
	deployCmd.Flags().BoolVarP(&fRun, "run", "", false, "Creates your app environment w/o webs or workers")
}

// nanoDeploy
func nanoDeploy(ccmd *cobra.Command, args []string) {

	// PreRun: bootVM

	//
	v := url.Values{}
	v.Add("reset", strconv.FormatBool(fForce))
	v.Add("run", strconv.FormatBool(fRun))

	//
	deploy := util.Sync{
		Model:   "deploy",
		Path:    fmt.Sprintf("%s/deploys?%v", config.ServerURL, v.Encode()),
		Verbose: fVerbose,
	}

	//
	deploy.Run(args)

	//
	switch deploy.Status {

	// for each successful deploy create/update the .nanobox/apps/<app>/.deployed
	// file
	case "complete":
		config.VMfile.DeployedIs(true)
		break

	// if a deploy ever errors, remove the deployed file; don't need to handle
	// an error here because it just means the file already doesn't exist
	case "errored":
		fmt.Println(`
! AN ERROR PREVENTED NANOBOX FROM BUILDING YOUR ENVIRONMENT !
- View the output above to diagnose the source of the problem
- You can also retry with --verbose for more detailed output
`)
		// this could probably be better
		if fWatch {
			config.VMfile.Mode = "background"
		}

		// deploy failed
		config.VMfile.DeployedIs(false)

		//
		nanoboxDown(nil, args)
		os.Exit(1)
	}
}
