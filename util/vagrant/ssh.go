// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

//
package vagrant

import (
	"os"
	"os/exec"

	"github.com/nanobox-io/nanobox/config"
)

// SSH is run manually (vs Run) because the output needs to be hooked up differntly
func SSH() error {

	//
	setContext(config.AppDir)

	cmd := exec.Command("vagrant", "ssh")

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// start the command; we need this to 'fire and forget' so that we can manually
	// capture and modify the commands output
	return cmd.Run()
}
