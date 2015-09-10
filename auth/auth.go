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

	api "github.com/pagodabox/nanobox-api-client"
	"github.com/pagodabox/nanobox-cli/config"
	"github.com/pagodabox/nanobox-cli/util"
	"github.com/pagodabox/nanobox-golang-stylish"
)

//
var (
	creds    *credentials //
	authfile string       //
)

// credentials represents all available/expected .authfile configurable options
type credentials struct {
	userslug  string `json:"user_slug"`  //
	authtoken string `json:"auth_token"` //
}

// init
func init() {

	// check for a ~/.nanobox/.auth file and create one if it's not found
	authfile = filepath.Clean(config.Root + "/.auth")
	if fi, _ := os.Stat(authfile); fi == nil {
		fmt.Printf(stylish.Bullet("Creating %s directory", authfile))
		if _, err := os.Create(authfile); err != nil {
			panic(err)
		}
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
	if creds.userslug == "" || creds.authtoken == "" {
		return false
	}

	// do a quick check to see if the cli needs to reauthenticate due to a user
	// changing their authenticate token via the dashboard.
	// if _, err := api.GetUser(creds.userslug); err != nil {
	// 	config.Log.Warn("Failed login attempt (%v): Credentials do not match! Reauthenticating...", err)
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

		userslug := util.Prompt("Username: ")
		password := util.PPrompt("Password: ")

		// authenticate
		return authenticate(userslug, password)
	}

	return creds.userslug, creds.authtoken
}

// Reauthenticate
func Reauthenticate() (string, string) {
	fmt.Println(`
It appears the Username or API token the CLI is trying to use does not match what
we have on record. To continue, please login to verify your account:
  `)

	userslug := util.Prompt("Username: ")
	password := util.PPrompt("Password: ")

	// authenticate
	return authenticate(userslug, password)
}

// authenticate
func authenticate(userslug, password string) (string, string) {

	fmt.Printf("\nAttempting login for %v... ", userslug)

	// get auth_token
	user, err := api.GetAuthToken(userslug, password)
	if err != nil {
		util.CPrint("[red]failure![reset]")
		fmt.Println("Unable to login... please verify your username and password are correct.")
		os.Exit(1)
	}

	//
	if err := saveCredentials(user.ID, user.AuthenticationToken); err != nil {
		util.LogFatal("[auth/auth] saveCredentials failed", err)
	}

	//
	util.CPrint("[green]success![reset]")

	return user.ID, user.AuthenticationToken
}

// writes user_slug and auth_token to .auth file
func saveCredentials(userid, authtoken string) error {

	//
	creds.userslug = userid
	creds.authtoken = authtoken

	//
	return ioutil.WriteFile(authfile, []byte(fmt.Sprintf("user_slug: %v\nauth_token: %v", userid, authtoken)), 0755)
}
