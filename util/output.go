// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

package util

import "fmt"

// Printv (print verbose) only prints a message if the 'verbose' flag is passed
func Printv(msg string, verbose bool) {
	if verbose {
		fmt.Printf(msg)
	}
}

//
func Debug(msg string, debug bool) {
	if debug {
		fmt.Printf(msg)
	}
}
