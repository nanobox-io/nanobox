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

// create default nanofilconfig
func init() {
	Nanofile = &NanofileConfig{
		CPUCap:   50,
		CPUs:     2,
		Domain:   fmt.Sprintf("%v.nano.dev", App),
		IP:       appNameToIP(App),
		Provider: "virtualbox",
		RAM:      1024,
	}
}

// Parse
func (c *NanofileConfig) Parse() error {
	fmt.Printf(stylish.Bullet("Parsing .nanofile"))

	//
	path := "./.nanofile"

	// look for a local .nanofile first...
	fmt.Printf(stylish.SubBullet("Searching for local .nanofile..."))
	if fi, _ := os.Stat(path); fi != nil {
		return c.parse(path)
	}

	path = NanoDir + "/.nanofile"

	// then look for a global .nanofile in the ~/.nanobox directory...
	fmt.Printf(stylish.SubBullet("Searching for global .nanofile..."))
	if fi, _ := os.Stat(path); fi != nil {
		return c.parse(path)
	}

	//
	fmt.Printf(stylish.SubBullet("- Using default configuration..."))
	fmt.Printf(stylish.Success())
	return nil
}

// parseNanofile
func (c *NanofileConfig) parse(path string) error {

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
		return fmt.Errorf("Nanobox failed to parse your .nanofile found at %s. Please ensure it is valid YAML and try again.", path)
	}

	fmt.Printf(stylish.Success())

	return nil
}
