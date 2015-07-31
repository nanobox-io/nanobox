// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/ghodss/yaml"

	"github.com/pagodabox/nanobox-golang-stylish"
)

type (
	// BoxfileConfig represents all available/expected Boxfile configurable options
	BoxfileConfig struct {
		Build Build //
	}

	//
	Build struct {
		Engine string `json:"engine"` //
	}
)

// Parse
func (bc *BoxfileConfig) Parse() error {
	fmt.Printf(stylish.Bullet("Parsing Boxfile"))

	//
	path := "./Boxfile"

	// look for a local Boxfile first...
	if fi, _ := os.Stat(path); fi != nil {
		return parseBoxfile(path, bc)
	}

	//
	fmt.Printf(stylish.SubBullet("- Using default configuration..."))
	fmt.Printf(stylish.Success())
	return nil
}

// parseBoxfile
func parseBoxfile(file string, bc *BoxfileConfig) error {

	fmt.Printf(stylish.SubBullet("- Configuring..."))

	fp, err := filepath.Abs(file)
	if err != nil {
		return err
	}

	//
	f, err := ioutil.ReadFile(fp)
	if err != nil {
		return err
	}

	//
	if err := yaml.Unmarshal(f, bc); err != nil {
		return fmt.Errorf("Nanobox failed to parse your Boxfile. Please ensure it is valid YAML and try again.")
	}

	fmt.Printf(stylish.Success())

	return nil
}
