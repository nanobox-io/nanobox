// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

//
package server

import ()

// Build issues a build to nanobox server
func Build(params string) error {

	res, err := Post("/builds?"+params, "text/plain", nil)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return nil
}
