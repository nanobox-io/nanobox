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
	"github.com/nanobox-io/nanobox/config"
)

// DeployUpdates
func DeployUpdates(status string) (err error) {

	switch status {

	// complete (deploy succeeded)
	case "complete":
		config.VMfile.DeployedIs(true)

	// errored (deploy failed)
	case "errored":
		config.VMfile.DeployedIs(false)
		config.VMfile.SuspendableIs(false)

		err = fmt.Errorf(`
! AN ERROR PREVENTED NANOBOX FROM BUILDING YOUR ENVIRONMENT !
- View the output above to diagnose the source of the problem
- You can also retry with --verbose for more detailed output
`)
	}

	return
}

// BuildUpdates receives a status update from mist.go and determines what
// to do based on the status. By default it will return, indicating to mist to
// stop listening.
func BuildUpdates(status string) (err error) {

	switch status {

	// complete (built succeeded)
	case "complete":

	// unavailable (deploy required)
	case "unavailable":
		err = fmt.Errorf("Before you can run a build, you must first deploy.")

	// errored (build failed)
	case "errored":
		config.VMfile.SuspendableIs(false)

		err = fmt.Errorf(`
! AN ERROR PREVENTED NANOBOX FROM BUILDING YOUR ENVIRONMENT !
- View the output above to diagnose the source of the problem
- You can also retry with --verbose for more detailed output
`)
	}

	return
}

// BootstrapUpdates
func BootstrapUpdates(status string) (err error) {

	switch status {

	// complete (bootstrap succeeded)
	case "complete":

	// errored (bootstrap failed)
	case "errored":
		config.VMfile.SuspendableIs(false)
		err = fmt.Errorf(`
! AN ERROR PREVENTED NANOBOX FROM BUILDING YOUR ENVIRONMENT !
- View the output above to diagnose the source of the problem
- You can also retry with --verbose for more detailed output
`)
	}

	return
}

// ImageUpdates
func ImageUpdates(status string) (err error) {

	switch status {

	// compelte (image update succeeded)
	case "complete":

	// errored (image update failed)
	case "errored":
		err = fmt.Errorf("Nanobox failed to update docker images")
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
