// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

package commands

import (
	"fmt"
	"io"
	"os"
	"regexp"

	api "github.com/pagodabox/nanobox-api-client"
	"github.com/pagodabox/nanobox-cli/auth"
	"github.com/pagodabox/nanobox-cli/config"
	"github.com/pagodabox/nanobox-cli/ui"
)

type (

	// FetchCommand satisfies the Command interface for listing a user's apps
	FetchCommand struct{}
)

// Help prints detailed help text for the app list command
func (c *FetchCommand) Help() {
	ui.CPrint(`
Description:
  Fetch an engine from nanobox.io

Usage:
  nanobox fetch ruby
  nanobox fetch nanobox/ruby
  nanobox fetch ruby-0.0.1
  nanobox fetch nanobox/ruby-0.0.1
  `)
}

// Run displays select information about all of a user's apps
func (c *FetchCommand) Run(opts []string) {

	// check for auth
	if !auth.IsAuthenticated() {
		fmt.Println("Before using the Pagoda Box CLI on this machine, please login to your account:")

		userslug := ui.Prompt("Username: ")
		password := ui.PPrompt("Password: ")

		// authenticate
		if err := auth.Authenticate(userslug, password); err != nil {
			ui.LogFatal("[main] auth.Authenticate() failed", err)
		}

		fmt.Println("To begin using the Pagoda Box CLI type 'pagoda' to see a list of commands.")
		os.Exit(0)
	}

	// pull the users api credentials
	creds, err := auth.Credentials()
	if err != nil {
		ui.LogFatal("[main] auth.Credentials() failed", err)
	}

	api.UserSlug = creds["user_slug"]
	api.AuthToken = creds["auth_token"]

	// if the credentials are empty attempt to reauthenticate
	if api.UserSlug == "" || api.AuthToken == "" {
		config.Console.Warn("No login credentials found! Reauthenticating...")
		auth.ReAuthenticate()
	}

	if len(opts) < 1 {
		config.Console.Fatal("Please provide the name of an engine you would like to fetch, (run 'nanobox fetch -h' for details)")
		os.Exit(1)
	}

	release := opts[0]

	// reExplodeRelease := regexp.MustCompile(`^(\w*\/)?(\w*)-?([0-9.]+)?$`)
	// release := reExplodeRelease.FindStringSubmatch(opts[0])
	// user, engine, version := release[0], release[1], release[2]

	// this should be able to explode most engine/version combinations
	reExplodeRelease := regexp.MustCompile(`^(.*[^\/])?-([version]*\d.+-?\w*)$`)
	match := reExplodeRelease.FindStringSubmatch(release)

	//
	engine := release
	if len(match) >= 1 {
		engine = match[1]
	}

	//
	e, err := api.GetEngine(api.UserSlug, engine)
	if err != nil {
		fmt.Println("ERR!!", err)
		config.Console.Info("No engines found on nanobox by the name '%v'", engine)
	}

	//
	version := e.ActiveReleaseID
	if len(match) >= 2 {
		version = match[2]
	}

	// pull directly from warehouse
	path := fmt.Sprintf("http://warehouse.nanobox.io/objects/releases/%v", version)

	// pull from odin
	// path := fmt.Sprintf("http://api.nanobox.io/v1/engines/%v/%v/releases/%v/download", api.UserSlug, e.Name, version)

	fmt.Println("FETCH!", path)

	//
	out, err := os.Create("release.tgz")
	if err != nil {
		ui.LogFatal("[commands.publish] os.Create() failed", err)
	}
	defer out.Close()

	//
	headers := map[string]string{
		"Userid":      e.WarehouseUser,
		"Key":         e.WarehouseKey,
		"Bucketid":    e.ID,
		"Objectalias": "releases/" + version,
	}

	fmt.Printf("HEADERS! %#v\n", headers)

	//
	req, err := api.NewRequest("GET", path, nil, headers)
	if err != nil {
		ui.LogFatal("[commands.publish] api.DoRawRequest() failed", err)
	}

	// req, err := api.NewRequest("GET", path, nil, nil)
	// if err != nil {
	// 	ui.LogFatal("[commands.publish] api.DoRawRequest() failed", err)
	// }

	//
	res, err := api.HTTPClient.Do(req)
	if err != nil {
		ui.LogFatal("[commands.publish] api.HTTPClient.Do() failed", err)
	}
	defer res.Body.Close()

	fmt.Println("RESPONSE!!", res)

	//
	if _, err := io.Copy(out, res.Body); err != nil {
		ui.LogFatal("[commands.publish] io.Copy() failed", err)
	}
}
