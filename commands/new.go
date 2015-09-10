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
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"

	"github.com/pagodabox/nanobox-cli/util"
	"github.com/pagodabox/nanobox-golang-stylish"
)

//
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Generates a new engine inside the current working directory",
	Long: `
Description:
  Generates a new engine inside the current working directory`,

	Run: nanoNew,
}

// nanoNew
func nanoNew(ccmd *cobra.Command, args []string) {

	if len(args) < 1 {
		fmt.Println(stylish.Error("app name required", "Please provide a name when generating a new engine, (run 'nanobox new -h' for details)"))
		os.Exit(1)
	}

	name := fmt.Sprintf("nanobox-%s", args[0])
	version := "0.0.1"

	fmt.Println(stylish.Header("initializing new app: %s", name))

	// create a new project by the name, unless it already exists
	if di, _ := os.Stat(name); di == nil {

		stylish.Bullet("Creating engine at")

		if err := os.Mkdir(name, 0755); err != nil {
			util.LogFatal("[commands/new] os.Mkdir() failed", err)
		}

		entry := fmt.Sprintf(`
name: %-18s     # the name of your project (required)
version: %-18s  # the current version of the project (required)
language:                    # the lanauge (ruby, golang, etc.) of the engine (required)
summary:                     # a 140 character short summary of the project (required)
stability: alpha             # the current stability of the project (alpha/beta/stable)
generic: false               # is the engine generic enough to encompass multiple frameworks
                             # within the given language
license: MIT                 # the license to be applied to this software

# a list of authors/contributors
authors:
  -
`, name, version)

		path := name + "/Enginefile"

		if _, err := os.Create(path); err != nil {
			util.LogFatal("[commands/new] os.Create() failed", err)
		}

		// write the Enginefile
		if err := ioutil.WriteFile(path, []byte(entry), 0644); err != nil {
			util.LogFatal("[commands/new] ioutil.WriteFile() failed", err)
		}

	} else {
		fmt.Printf("A project by the name '%s' already exists at this location...\n", name)
	}

	stylish.Bullet("Default Enginefile created at")
}
