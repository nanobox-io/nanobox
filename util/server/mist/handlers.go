// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

//
package mist

import (
	"fmt"

	"github.com/nanobox-io/nanobox-cli/config"
	"github.com/nanobox-io/nanobox-golang-stylish"
)

// HandleDeployStream
func HandleDeployStream(status string) (listen bool) {

	switch status {

	// continue listening until one of the following statuses is received
	default:
		listen = true

	// complete (deploy succeeded)
	case "complete":
		config.VMfile.DeployedIs(true)

	// errored (deploy failed)
	case "errored":
		config.VMfile.DeployedIs(false)

		fmt.Println(`
FAIL:
   ! AN ERROR PREVENTED NANOBOX FROM BUILDING YOUR ENVIRONMENT !
   - View the output above to diagnose the source of the problem
   - You can also retry with --verbose for more detailed output`)
	}

	return
}

// HandleBuildStream receives a status update from mist.go and determines what
// to do based on the status. By default it will return, indicating to mist to
// stop listening.
func HandleBuildStream(status string) (listen bool) {

	switch status {

	// continue listening until one of the following statuses is received
	default:
		listen = true

	// complete (built succeeded)
	case "complete":

	// unavailable (deploy required)
	case "unavailable":
		fmt.Printf(stylish.ErrBullet("Before you can run a build, you must first deploy."))

	// errored (build failed)
	case "errored":

		fmt.Println(`
FAIL:
   ! AN ERROR PREVENTED NANOBOX FROM BUILDING YOUR ENVIRONMENT !
   - View the output above to diagnose the source of the problem
   - You can also retry with --verbose for more detailed output`)
	}

	return
}

// HandleBootstrapStream
func HandleBootstrapStream(status string) (listen bool) {

	switch status {

	// continue listening until one of the following statuses is received
	default:
		listen = true

	// complete (bootstrap succeeded)
	case "complete":

	// errored (bootstrap failed)
	case "errored":
	}

	return
}

// HandleUpdateStream
func HandleUpdateStream(status string) (listen bool) {

	switch status {

	// continue listening until one of the following statuses is received
	default:
		listen = true

	// compelte (image update succeeded)
	case "complete":

	// errored (image update failed)
	case "errored":
		fmt.Printf("Nanobox failed to update docker images")
	}

	return
}

// PrintLogStream prints a message received as is
func PrintLogStream(log Log) {
	fmt.Printf(log.Content)
}

// ProcessLogStream processes a log before printing it
func ProcessLogStream(log Log) {

	// if the CLI is silenced don't process any logs
	if !config.Silence {
		ProcessLog(log)
	}
}
