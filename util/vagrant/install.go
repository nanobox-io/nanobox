// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

//
package vagrant

// Install downloads the nanobox vagrant and adds it to the list of vagrant boxes
func Install() error {

	// download nanobox image
	if err := download(); err != nil {
		return err
	}

	// add nanobox image
	return add()
}
