// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

package config

import (
	"fmt"
	"os"
)

// ParseBoxfile
func ParseBoxfile() *BoxfileConfig {

	boxfile := &BoxfileConfig{path: "./Boxfile"}

	//
	if _, err := os.Stat(boxfile.path); err != nil {
		fmt.Println("Boxfile not found in current working directory, unable to continue...")
		os.Exit(1)
	}

	//
	if err := ParseConfig(boxfile.path, boxfile); err != nil {
		fmt.Printf("Nanobox failed to parse your Boxfile. Please ensure it is valid YAML and try again.\n")
		os.Exit(1)
	}

	//
	return boxfile
}
