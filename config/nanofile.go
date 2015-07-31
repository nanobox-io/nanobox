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

// NanofileConfig represents all available/expected Boxfile configurable options
type NanofileConfig struct {
	CPUCap   int    `json:"cpu_cap"`  //
	CPUs     int    `json:"cpus"`     //
	Domain   string `json:"domain"`   //
	IP       string `json:"ip"`       //
	Provider string `json:"provider"` //
	RAM      int    `json:"ram"`      //
}

// Parse
func (nc *NanofileConfig) Parse() error {
	fmt.Printf(stylish.Bullet("Parsing .nanofile"))

	//
	path := "./.nanofile"

	// look for a local .nanofile first...
	fmt.Printf(stylish.SubBullet("Searching for local .nanofile..."))
	if fi, _ := os.Stat(path); fi != nil {
		return parseNanofile(path, nc)
	}

	path = NanoDir + "/.nanofile"

	// then look for a global .nanofile in the ~/.nanobox directory...
	fmt.Printf(stylish.SubBullet("Searching for global .nanofile..."))
	if fi, _ := os.Stat(path); fi != nil {
		return parseNanofile(path, nc)
	}

	//
	fmt.Printf(stylish.SubBullet("- Using default configuration..."))
	fmt.Printf(stylish.Success())
	return nil
}

// parseNanofile
func parseNanofile(file string, nc *NanofileConfig) error {

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
	if err := yaml.Unmarshal(f, nc); err != nil {
		return fmt.Errorf("Nanobox failed to parse your .nanofile found at %s. Please ensure it is valid YAML and try again.", file)
	}

	fmt.Printf(stylish.Success())

	return nil
}
