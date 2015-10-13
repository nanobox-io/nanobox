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

// DeployUpdates
func DeployUpdates(status string) (listen bool) {

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
! AN ERROR PREVENTED NANOBOX FROM BUILDING YOUR ENVIRONMENT !
- View the output above to diagnose the source of the problem
- You can also retry with --verbose for more detailed output
`)

		config.VMfile.SuspendableIs(false)
	}

	return
}

// BuildUpdates receives a status update from mist.go and determines what
// to do based on the status. By default it will return, indicating to mist to
// stop listening.
func BuildUpdates(status string) (listen bool) {

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
! AN ERROR PREVENTED NANOBOX FROM BUILDING YOUR ENVIRONMENT !
- View the output above to diagnose the source of the problem
- You can also retry with --verbose for more detailed output
`)

		config.VMfile.SuspendableIs(false)
	}

	return
}

// BootstrapUpdates
func BootstrapUpdates(status string) (listen bool) {

	switch status {

	// continue listening until one of the following statuses is received
	default:
		listen = true

	// complete (bootstrap succeeded)
	case "complete":

	// errored (bootstrap failed)
	case "errored":
		config.VMfile.SuspendableIs(false)
	}

	return
}

// ImageUpdates
func ImageUpdates(status string) (listen bool) {

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

// ProcessLogStream processes a log before printing it; if the CLI is silenced
// don't process any logs
func ProcessLogStream(log Log) {
	if !config.Silent {
		ProcessLog(log)
	}
}
