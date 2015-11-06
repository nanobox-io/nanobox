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
	"github.com/nanobox-io/nanobox/config"
	fileutil "github.com/nanobox-io/nanobox/util/file"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// ExtractArchive splits args on "/" looking for a user and archive:
// - user/engine-name
// - user/engine-name=0.0.1
func ExtractArchive(s string) (user, archive string) {

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

// ExtractEngine splits on the archive to find the engine and the release (version)
func ExtractEngine(archive string) (engine, version string) {

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

// GetEngine
func GetEngine(user, archive, version string) (*http.Response, error) {

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

//
func MountLocal() {

	//
	boxfile := config.ParseBoxfile()

	// if an custom engine path is provided, add it to the synced_folders
	if engine := boxfile.Build.Engine; engine != "" {
		if _, err := os.Stat(engine); err == nil {

			//
			mountDir := filepath.Join(config.AppDir, filepath.Base(engine))
			if _, err := os.Stat(mountDir); err != nil {
				if err := os.Mkdir(mountDir, 0755); err != nil {
					config.Fatal("[commands/init] os.Mkdir() failed", err.Error())
				}
			}

			//
			enginefile := filepath.Join(engine, "./Enginefile")

			// if no engine file is found just return
			if _, err := os.Stat(enginefile); err != nil {
				fmt.Printf("No enginefile found at '%v', Exiting...\n", engine)
				os.Exit(1)
			}

			//
			mount := &struct {
				Overlays []string `json:"overlays"`
			}{}

			// parse the ./Enginefile into the new mount
			if err := config.ParseConfig(enginefile, mount); err != nil {
				fmt.Printf("Nanobox failed to parse your Enginefile. Please ensure it is valid YAML and try again.\n")
				os.Exit(1)
			}

			// iterate through each overlay fetching it and adding it to the list of 'files'
			// to be tarballed
			for _, overlay := range mount.Overlays {

				// extract a user and archive (desired engine) from args[0]
				user, archive := ExtractArchive(overlay)

				// extract an engine and version from the archive
				e, version := ExtractEngine(archive)

				//
				res, err := GetEngine(user, e, version)
				if err != nil {
					config.Fatal("[util/engine/engine] http.Get() failed", err.Error())
				}
				defer res.Body.Close()

				//
				switch res.StatusCode / 100 {
				case 2, 3:
					break
				case 4, 5:
					os.Stderr.WriteString(stylish.ErrBullet("Unable to fetch '%v' overlay, exiting...", e))
					os.Exit(1)
				}

				//
				if err := fileutil.Untar(mountDir, res.Body); err != nil {
					config.Fatal("[util/engine/engine] file.Untar() failed", err.Error())
				}
			}

			abs, err := filepath.Abs(engine)
			if err != nil {
				config.Fatal("[util/engine/engine] filepath.Abs() failed", err.Error())
			}

			// pull the remainin engine files over
			for _, f := range []string{"bin", "Enginefile", "lib", "templates", "files"} {

				path := filepath.Join(abs, f)

				// just skip any files that aren't found; any required files will be
				// caught before publishing, here it doesn't matter
				if _, err := os.Stat(path); err != nil {
					continue
				}

				// not handling error here because an error simply means the file doesn't
				// exist and therefor wont be copied
				if err := fileutil.Copy(path, mountDir); err != nil {
					config.Fatal("[util/engine/engine] file.Copy() failed", err.Error())
				}
			}
		}
	}
}
