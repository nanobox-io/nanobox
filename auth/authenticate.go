// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

package auth

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	api "github.com/pagodabox/nanobox-api-client"
	"github.com/pagodabox/nanobox-cli/config"
	"github.com/pagodabox/nanobox-cli/ui"
)

// IsAuthenticated checks to see if there is a .auth file in the home dir
func IsAuthenticated() bool {
	f, _ := os.Stat(config.AuthFile)
	return (f != nil)
}

// Authenticate prompts a user for their username and password to authenticate
// the CLI with the Pagoda Box API. If able to authenticate it continues to
// 'install' the CLI by creating the necessary folders and files for running the
// CLI. It takes the token returned in the authentication process and saves that
// to a .token file to be used with each CLI command's resulting API request.
func Authenticate(userslug, password string) error {

	fmt.Printf("\nAttempting login for %v... ", userslug)

	// get auth_token
	user, err := api.GetAuthToken(userslug, password)
	if err != nil {
		ui.CPrint("[red]failure![reset]")
		fmt.Println("Unable to login... please verify your username and password are correct.")
		return err
	}

	// write user_slug and auth_token to .auth file
	if err := ioutil.WriteFile(config.AuthFile, []byte("user_slug "+user.ID+"\nauth_token "+user.AuthenticationToken), 0755); err != nil {
		return err
	}

	//
	ui.CPrint("[green]success![reset]")

	return nil
}

// ReAuthenticate is run after a 'ping' done to the API that returns the error
// "Invalid authentication token". This indicates that token stored in the .token
// file does not match what the API expects. Authenticate called w/o a username
// or password causing a prompt for both. This will cause the .token file to store
// the 'new' auth token from the login attempt
func ReAuthenticate() {
	fmt.Println(`
It appears the Username or API token the CLI is trying to use does not match what
we have on record. To continue, please login to verify your account:
  `)

	userslug := ui.Prompt("Username: ")
	password := ui.PPrompt("Password: ")

	// authenticate
	if err := Authenticate(userslug, password); err != nil {
		ui.LogFatal("[helpers.authenticate] Authenticate() failed", err)
	}

	fmt.Println("You may now continue using the nanobox CLI.")
	os.Exit(0)
}

// Credentials reads the .auth file and returns the API authentication options
func Credentials() (map[string]string, error) {

	// attempt to open file
	f, err := os.Open(config.AuthFile)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	opts := make(map[string]string)
	scanner := bufio.NewScanner(f)

	// Read line by line
	for scanner.Scan() {

		// extract key/value pair
		fields := strings.Fields(scanner.Text())

		// insert key/value pair into map
		opts[fields[0]] = fields[1]
	}

	return opts, nil
}
