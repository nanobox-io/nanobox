// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.

//
package server

import ()

// IsContainerExec
func IsContainerExec(args []string) (found bool) {

	// fetch services to see if the command is trying to run on a specific container
	var services []Service
	res, err := Get("/services", &services)
	if err != nil {
		Fatal("[util/server/exec] Get() failed - ", err.Error())
	}
	defer res.Body.Close()

	//
	for _, service := range services {

		// range over the services to find a potential match for args[0] or make an
		// exception for 'build1' since that wont show up on the list.
		if args[0] == service.Name || args[0] == "build1" {
			found = true
		}
	}

	return
}
