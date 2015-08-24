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
	"path/filepath"

	"github.com/ghodss/yaml"

	"github.com/pagodabox/nanobox-golang-stylish"
)

// AuthfileConfig represents all available/expected Enginefile configurable options
type AuthfileConfig struct {
	UserSlug  string `json:"user_slug"`  //
	AuthToken string `json:"auth_token"` //
}

// Parse
func (c *AuthfileConfig) Parse() error {
	fmt.Printf(stylish.Bullet("Parsing Authfile..."))
	return c.parse(AuthFile)
}

// create default authfile config
func init() {
	Authfile = &AuthfileConfig{}
}

// parse
func (c *AuthfileConfig) parse(path string) error {

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
		return fmt.Errorf("Nanobox failed to parse your Authfile. Please ensure it is valid YAML and try again.")
	}

	fmt.Printf(stylish.Success())

	return nil
}
