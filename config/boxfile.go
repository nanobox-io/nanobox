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

// create default boxfile config
func init() {
	Boxfile = &BoxfileConfig{}
}

// Parse
func (c *BoxfileConfig) Parse() error {
	fmt.Printf(stylish.Bullet("Parsing Boxfile"))

	//
	path := "./Boxfile"

	// look for a Boxfile...
	if fi, _ := os.Stat(path); fi != nil {
		return c.parse(path)
	}

	//
	fmt.Printf(stylish.SubBullet("- Using default configuration..."))
	fmt.Printf(stylish.Success())
	return nil
}

// parse
func (c *BoxfileConfig) parse(path string) error {

	fmt.Printf(stylish.SubBullet("- Configuring..."))

	fp, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	//
	f, err := ioutil.ReadFile(fp)
	if err != nil {
		return err
	}

	//
	if err := yaml.Unmarshal(f, c); err != nil {
		return fmt.Errorf("Nanobox failed to parse your Boxfile. Please ensure it is valid YAML and try again.")
	}

	fmt.Printf(stylish.Success())

	return nil
}
