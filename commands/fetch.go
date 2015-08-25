// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

package commands

//
import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/cobra"

	api "github.com/pagodabox/nanobox-api-client"
	// "github.com/pagodabox/nanobox-cli/auth"
	"github.com/pagodabox/nanobox-cli/config"
	"github.com/pagodabox/nanobox-cli/ui"
	"github.com/pagodabox/nanobox-golang-stylish"
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
  - user/engine-name
  - engine-name-0.0.1
  - user/engine-name-0.0.1
	`,

	Run: nanoFetch,
}

//
func init() {

	// no default is set here because we define the value later, once we know the
	// name and version of the engine they are fetching
	fetchCmd.Flags().StringVarP(&fFile, "ouput-document", "O", "", "specify where to save the engine")
}

// nanoFetch
func nanoFetch(ccmd *cobra.Command, args []string) {

	//
	// api.UserSlug, api.AuthToken = auth.Authenticate()

	if len(args) == 0 {
		config.Console.Fatal("Please provide the name of an engine you would like to fetch, (run 'nanobox fetch -h' for details)")
		os.Exit(1)
	}

	fmt.Printf(stylish.Bullet(fmt.Sprintf("Attempting to fetch '%v'", args[0])))

	//
	var archive, engine, user, version string // various string values used to store pieces of the engine
	var split []string                        // used in strings.Split()
	var dest io.Writer                        // the destination used in io.Copy()

	// split args on "/" looking for a user:
	// user/engine-name
	// user/engine-name-0.0.1
	split = strings.Split(args[0], "/")

	// switch on the length to determin if the split resulted in a user and a engine
	// or just an engine
	switch len(split) {

	// if len is 1 then only a download was found (no user specified)
	case 1:
		archive = split[0]

		// if len is 2 then a user was found (from which to pull the download)
	case 2:
		user = split[0]
		archive = split[1]

		// if some other length was found, then the split 'failed' (meaning the
		// format of the fetch was probably incorrect)
	default:
		fmt.Printf("%v is not a valid format when fetching an engine (see help).\n", args[0])
		os.Exit(1)
	}

	// split on the archive to find the engine and the release (the release version
	// being the last index in the split)
	split = strings.Split(archive, "-")

	// the engine name is reconstructed from the previous split, joining each index
	// from the split upto but not including the final element (which should be the
	// version)
	engine = strings.Join(split[0:(len(split)-1)], "-")

	// the version is extracted from the split (being the last index in the split)
	version = split[len(split)-1]

	//
	if _, err := api.GetEngine(api.UserSlug, engine); err != nil {
		config.Console.Info("Failed to GET engine '%v' - %v", engine, err)
		os.Exit(1)
	}

	//
	path := fmt.Sprintf("http://api.nanobox.io/v1/engines/%v/releases/%v/download", engine, version)
	if user != "" {
		path = fmt.Sprintf("http://api.nanobox.io/v1/engines/%v/%v/releases/%v/download", user, engine, version)
	}

	fmt.Printf(stylish.Bullet(fmt.Sprintf("Fetching engine at '%v'", path)))

	res, err := http.Get(path)
	if err != nil {
		ui.LogFatal("[commands.fetch] http.Get() failed", err)
	}
	defer res.Body.Close()

	//
	switch {

	// write the download to the local file system
	case fFile != "":
		fmt.Printf(stylish.Bullet(fmt.Sprintf("Saving engine as '%v'", fFile)))

		//
		release, err := os.Create(fFile)
		if err != nil {
			ui.LogFatal("[commands.fetch] os.Create() failed", err)
		}
		defer release.Close()

		//
		dest = release

		// pipe the ouput to os.Stdout
	default:
		fmt.Printf(stylish.Bullet("Piping download to stdout"))
		dest = os.Stdout
	}

	// write the file
	if _, err := io.Copy(dest, res.Body); err != nil {
		ui.LogFatal("[commands.fetch] io.Copy() failed", err)
	}
}
