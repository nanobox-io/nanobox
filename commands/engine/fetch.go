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
	"github.com/nanobox-io/nanobox-golang-stylish"
	// "github.com/nanobox-io/nanobox/auth"
	engineutil "github.com/nanobox-io/nanobox/util/engine"
	"github.com/spf13/cobra"
	"io"
	"os"
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
	user, archive := engineutil.ExtractArchive(args[0])

	// extract an engine and version from the archive
	engine, version := engineutil.ExtractEngine(archive)

	// pull the engine from nanobox.io
	res, err := engineutil.GetEngine(user, engine, version)
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

	// determine if destination will be a file or stdout (stdout by default)
	dest := os.Stdout
	defer dest.Close()

	// write the download to the local file system
	if fFile != "" {

		//
		f, err := os.Create(fFile)
		if err != nil {
			os.Stderr.WriteString(stylish.ErrBullet("Unable to save file, exiting... %v", err.Error()))
			return
		}

		// if the file was created successfully then set it as the destination
		os.Stderr.WriteString(stylish.Bullet("Saving engine as '%s'", fFile))
		dest = f
	}

	// write the file
	if _, err := io.Copy(dest, res.Body); err != nil {
		os.Stderr.WriteString(fmt.Sprintf("[commands.fetch] io.Copy() failed - %s", err.Error()))
	}
}
