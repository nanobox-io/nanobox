// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

package main

import (
	// api "github.com/pagodabox/nanobox-api-client"
	"github.com/pagodabox/nanobox-cli/commands"
)

// main
func main() {

	// check for updates
	// checkUpdate()

	// do a quick ping to make sure we can communicate properly with the API
	// if err := api.DoRawRequest(nil, "GET", "https://api.pagodabox.io/v1/ping", nil, nil); err != nil {
	// 	config.Fatal("[main] The CLI was unable to communicate with the API", err)
	// }

	//
	commands.NanoboxCmd.Execute()
}
