// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

//
package engine

import (
	"fmt"
	api "github.com/nanobox-io/nanobox-api-client"
	"github.com/nanobox-io/nanobox-golang-stylish"
	// "github.com/nanobox-io/nanobox/auth"
	"github.com/spf13/cobra"
	"io"
	"net/http"
	"os"
	"strings"
)

//
var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "Fetches an engine from nanobox.io",
	Long: `
Description:
  Fetches an engine from nanobox.io

  Allowed formats when fetching an engine
  - engine-name
	- engine-name=0.0.1
  - user/engine-name
  - user/engine-name=0.0.1
	`,

	Run: fetch,
}

//
func init() {

	// no default is set here because we define the value later, once we know the
	// name and version of the engine they are fetching
	fetchCmd.Flags().StringVarP(&fFile, "ouput-document", "O", "", "specify where to save the engine")
}

// fetch
func fetch(ccmd *cobra.Command, args []string) {

	//
	// api.UserSlug, api.AuthToken = auth.Authenticate()

	if len(args) == 0 {
		os.Stderr.WriteString("Please provide the name of an engine you would like to fetch, (run 'nanobox engine fetch -h' for details)")
		os.Exit(1)
	}

	os.Stderr.WriteString(stylish.Bullet("Attempting to fetch '%v'", args[0]))

	// extract a user and archive (desired engine) from args[0]
	user, archive := extractArchive(args[0])

	// extract an engine and version from the archive
	engine, version := extractEngine(archive)

	//
	res, err := getEngine(user, engine, version)
	if err != nil {
		Config.Fatal("[commands/engine/fetch] http.Get() failed", err.Error())
	}
	defer res.Body.Close()

	//
	switch res.StatusCode / 100 {
	case 2, 3:
		break
	case 4:
		os.Stderr.WriteString(stylish.ErrBullet("No release by that version found for engine '%v'", engine))
		os.Exit(1)
	case 5:
		os.Stderr.WriteString(stylish.ErrBullet("Failed to fetch release (%v).", res.Status))
		os.Exit(1)
	}

	// determine the destination where the release will end up (file or stdout)
	dest := setDestination(fFile)
	defer dest.Close()

	// write the file
	if _, err := io.Copy(dest, res.Body); err != nil {
		os.Stderr.WriteString(fmt.Sprintf("[commands.fetch] io.Copy() failed - %s", err.Error()))
	}
}

// extractArchive splits args on "/" looking for a user and archive:
// - user/engine-name
// - user/engine-name=0.0.1
func extractArchive(s string) (user, archive string) {

	split := strings.Split(s, "/")

	// switch on the length to determine if the split resulted in a user and a engine
	// or just an engine
	switch len(split) {

	// if len is 1 then only a download was found (no user specified)
	case 1:
		archive = split[0]

		// if len is 2 then a user was found (from which to pull the download)
	case 2:
		user = split[0]
		archive = split[1]

	// any other number or args
	default:
		// fmt.Printf("%v is not a valid format when fetching an engine (see help).\n", args[0])
		os.Exit(1)
	}

	return
}

// extractEngine splits on the archive to find the engine and the release (version)
func extractEngine(archive string) (engine, version string) {

	// split on '=' looking for a version
	split := strings.Split(archive, "=")

	// switch on the length to determine if the split resulted in a engine and version
	// or just an engine
	switch len(split) {

	// if len is 1 then just an engine was found (no version specified)
	case 1:
		engine = split[0]

	// if len is 2 then an engine and version were found
	case 2:
		engine = split[0]
		version = split[1]
	}

	return
}

//
func getEngine(user, archive, version string) (*http.Response, error) {

	//
	engine, err := api.GetEngine(user, archive)
	if err != nil {
		os.Stderr.WriteString(stylish.ErrBullet("No official engine, or engine for that user found."))
		return nil, err
	}

	// if no version is provided, fetch the latest release
	if version == "" {
		version = engine.ActiveReleaseID
	}

	//
	path := fmt.Sprintf("http://api.nanobox.io/v1/engines/%v/releases/%v/download", archive, version)

	// if a user is found, pull the engine from their engines
	if user != "" {
		path = fmt.Sprintf("http://api.nanobox.io/v1/engines/%v/%v/releases/%v/download", user, archive, version)
	}

	os.Stderr.WriteString(stylish.Bullet("Fetching engine at '%s'", path))

	//
	return http.Get(path)
}

// setDestination determines if the file is to be streamed to stdout or to a file
func setDestination(path string) (dest io.WriteCloser) {

	switch {

	// pipe the ouput to os.Stdout
	default:
		dest = os.Stdout

		// write the download to the local file system
	case path != "":
		os.Stderr.WriteString(stylish.Bullet("Saving engine as '%s'", path))

		var err error

		//
		if dest, err = os.Create(path); err != nil {
			os.Stderr.WriteString(stylish.ErrBullet(err.Error()))
			os.Exit(1)
		}
	}

	return
}
