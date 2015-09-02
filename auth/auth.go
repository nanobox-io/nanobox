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

	api "github.com/pagodabox/nanobox-api-client"
	"github.com/pagodabox/nanobox-cli/config"
	"github.com/pagodabox/nanobox-cli/util"
	"github.com/pagodabox/nanobox-golang-stylish"
)

// authenticated checks to see if there is a .auth file in the home dir
func authenticated() bool {

	//
	if err := config.Authfile.Parse(); err != nil {
		util.LogFatal("auth/auth] config.Authfile.Parse() failed", err)
	}

	//
	if config.Authfile.UserSlug == "" || config.Authfile.AuthToken == "" {
		return false
	}

	// do a quick check to see if the cli needs to reauthenticate due to a user
	// changing their authenticate token via the dashboard.
	// if _, err := api.GetUser(config.Authfile.UserSlug); err != nil {
	// 	config.Log.Warn("Failed login attempt (%v): Credentials do not match! Reauthenticating...", err)
	// 	Reauthenticate()
	// }

	return true
}

//
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

	return config.Authfile.UserSlug, config.Authfile.AuthToken
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
	config.Authfile.UserSlug = userid
	config.Authfile.AuthToken = authtoken

	//
	return ioutil.WriteFile(config.AuthFile, []byte(fmt.Sprintf("user_slug: %v\nauth_token: %v", userid, authtoken)), 0755)
}
