// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

//
package vagrant

import (
	"fmt"
	"github.com/nanobox-io/nanobox-golang-stylish"
	"os/exec"
)

// Suspend runs a "vagrant suspend"
func Suspend() error {

	// suspend the vm
	fmt.Printf("\n%s", stylish.Bullet("Suspending nanobox..."))
	if err := runInContext(exec.Command("vagrant", "suspend")); err != nil {
		return err
	}
	fmt.Printf(stylish.Bullet("Exiting"))

	return nil
}
