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

// EnginefileConfig represents all available/expected Enginefile configurable options
type EnginefileConfig struct {
	Authors   []string `json:"authors"`   //
	Generic   string   `json:"generic"`   //
	Language  string   `json:"language"`  //
	License   string   `json:"license"`   //
	Name      string   `json:"name"`      //
	Readme    string   `json:"readme"`    //
	Stability string   `json:"stability"` //
	Summary   string   `json:"summary"`   //
	Version   string   `json:"version"`   //
}

// create default enginefile config
func init() {
	Enginefile = &EnginefileConfig{}
}

// Parse
func (c *EnginefileConfig) Parse() error {
	fmt.Printf(stylish.Bullet("Parsing Enginefile"))

	//
	path := "./Enginefile"

	//
	if _, err := os.Stat(path); err != nil {
		fmt.Println("Enginefile not found. Be sure to publish from a project directory. Exiting... ")
		os.Exit(1)
	}

	return c.parse(path)
}

// parse
func (c *EnginefileConfig) parse(path string) error {

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
		return fmt.Errorf("Nanobox failed to parse your Enginefile. Please ensure it is valid YAML and try again.")
	}

	fmt.Printf(stylish.Success())

	return nil
}
