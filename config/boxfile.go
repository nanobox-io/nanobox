// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

//
package config

import (
	"fmt"
	"os"
)

// BoxfileConfig represents all available/expected Boxfile configurable options
type BoxfileConfig struct {
	Build struct {
		Engine string `json:"engine"`
	}
}

// ParseBoxfile
func ParseBoxfile() *BoxfileConfig {

	boxfile := &BoxfileConfig{}
	boxfilePath := "./Boxfile"

	//
	if _, err := os.Stat(boxfilePath); err != nil {
		return boxfile
	}

	//
	if err := ParseConfig(boxfilePath, boxfile); err != nil {
		fmt.Printf("Nanobox failed to parse Boxfile. Please ensure it is valid YAML and try again.\n")
		Exit(1)
	}

	//
	return boxfile
}
