// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

//
package server

import (
	"fmt"

	"github.com/nanobox-io/nanobox-cli/util/server/mist"
	"github.com/nanobox-io/nanobox-golang-stylish"
)

// Logs diplayes historical logs from the server
func Logs(params string) error {

	logs := []mist.Log{}

	//
	res, err := Get("/logs?"+params, &logs)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	//
	fmt.Printf(stylish.Bullet("Showing last %v entries:", len(logs)))
	for _, log := range logs {
		mist.ProcessLog(log)
	}

	return nil
}
