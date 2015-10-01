// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

package util

//
import (
	"fmt"
	"os"
	"os/exec"

	"github.com/nanobox-io/nanobox-cli/config"
	"github.com/nanobox-io/nanobox-golang-stylish"
)

// SudoExec runs a command as sudo
func SudoExec(command, msg string) {
	fmt.Printf(stylish.Bullet(msg))

	//
	cmd := exec.Command("/bin/sh", "-c", fmt.Sprintf("sudo %v %v", os.Args[0], command))

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// run command
	if err := cmd.Run(); err != nil {
		config.Fatal("[utils/exec]", err.Error())
	}
}
