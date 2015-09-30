// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

package auth

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	api "github.com/nanobox-io/nanobox-api-client"
	"github.com/nanobox-io/nanobox-cli/config"
	"github.com/nanobox-io/nanobox-cli/util"
	"github.com/nanobox-io/nanobox-golang-stylish"
)

//
var (
	creds    *credentials //
	authfile string       //
)

// credentials represents all available/expected .authfile configurable options
type credentials struct {
	Userslug  string `json:"user_slug"`  //
	Authtoken string `json:"auth_token"` //
}

// init
func init() {

	// check for a ~/.nanobox/.auth file and create one if it's not found
	authfile = filepath.Clean(config.Root + "/.auth")
	if _, err := os.Stat(authfile); err != nil {
		f, err := os.Create(authfile)
		if err != nil {
			config.Fatal("[auth/auth] os.Create() failed", err.Error())
		}
		defer f.Close()
	}

	creds = &credentials{}

	//
	if err := config.ParseConfig(authfile, creds); err != nil {
		fmt.Printf("Nanobox failed to parse the .auth file.\n")
		os.Exit(1)
	}
}

// authenticated checks to see if there is a .auth file in the home dir
func authenticated() bool {

	//
	if creds.Userslug == "" || creds.Authtoken == "" {
		return false
	}

	// do a quick check to see if the cli needs to reauthenticate due to a user
	// changing their authenticate token via the dashboard.
	// if _, err := api.GetUser(creds.Userslug); err != nil {
	// 	config.Log.Warn("Failed login attempt (%s): Credentials do not match! Reauthenticating...", err.Error())
	// 	Reauthenticate()
	// }

	return true
}

// Authenticate
func Authenticate() (string, string) {
	fmt.Printf(stylish.Bullet("Authenticating..."))

	//
	if !authenticated() {
		fmt.Println("Before continuing, please login to your account:")

		Userslug := util.Prompt("Username: ")
		password := util.PPrompt("Password: ")

		// authenticate
		return authenticate(Userslug, password)
	}

	return creds.Userslug, creds.Authtoken
}

// Reauthenticate
func Reauthenticate() (string, string) {
	fmt.Println(`
It appears the Username or API token the CLI is trying to use does not match what
we have on record. To continue, please login to verify your account:
  `)

	Userslug := util.Prompt("Username: ")
	password := util.PPrompt("Password: ")

	// authenticate
	return authenticate(Userslug, password)
}

// authenticate
func authenticate(Userslug, password string) (string, string) {

	fmt.Printf("\nAttempting login for %v... ", Userslug)

	// get auth_token
	user, err := api.GetAuthToken(Userslug, password)
	if err != nil {
		util.Printc("[red]failure![reset]")
		fmt.Println("Unable to login... please verify your username and password are correct.")
		os.Exit(1)
	}

	//
	if err := saveCredentials(user.ID, user.AuthenticationToken); err != nil {
		config.Fatal("[auth/auth] saveCredentials failed", err.Error())
	}

	//
	util.Printc("[green]success![reset]")

	return user.ID, user.AuthenticationToken
}

// writes user_slug and auth_token to .auth file
func saveCredentials(userid, authtoken string) error {

	//
	creds.Userslug = userid
	creds.Authtoken = authtoken

	//
	return ioutil.WriteFile(authfile, []byte(fmt.Sprintf("user_slug: %v\nauth_token: %v", userid, authtoken)), 0644)
}
