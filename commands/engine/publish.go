// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

//
package engine

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/spf13/cobra"

	api "github.com/nanobox-io/nanobox-api-client"
	"github.com/nanobox-io/nanobox-cli/auth"
	"github.com/nanobox-io/nanobox-cli/config"
	// "github.com/nanobox-io/nanobox-cli/util/file"
	"github.com/nanobox-io/nanobox-cli/util/s3"
	"github.com/nanobox-io/nanobox-golang-stylish"
)

var tw *tar.Writer

//
var publishCmd = &cobra.Command{
	Use:   "publish",
	Short: "Publishes an engine to nanobox.io",
	Long:  ``,

	Run: publish,
}

// publish
func publish(ccmd *cobra.Command, args []string) {
	//
	stylish.Header("publishing engine")

	//
	api.UserSlug, api.AuthToken = auth.Authenticate()

	// create a new release
	fmt.Printf(stylish.Bullet("Creating release..."))
	release := &api.EngineRelease{}

	//
	if _, err := os.Stat("./Enginefile"); err != nil {
		fmt.Println("Enginefile not found. Be sure to publish from a project directory. Exiting... ")
		os.Exit(1)
	}

	if err := config.ParseConfig("./Enginefile", release); err != nil {
		fmt.Printf("Nanobox failed to parse your Enginefile. Please ensure it is valid YAML and try again.\n")
		os.Exit(1)
	}

	fmt.Printf(stylish.Bullet("Verifying engine is publishable..."))

	// determine if any required fields (name, version, language, summary) are missing,
	// if any are found to be missing exit 1
	// NOTE: I do this using fallthrough for asthetics onlye. The message is generic
	// enough that all cases will return the same message, and this looks better than
	// a single giant case (var == "" || var == "" || ...)
	switch {
	case release.Language == "":
		fallthrough
	case release.Name == "":
		fallthrough
	case release.Summary == "":
		fallthrough
	case release.Version == "":
		fmt.Printf(stylish.Error("required fields missing", `Your Enginefile is missing one or more of the following required fields for publishing:

  name:      # the name of your project
  version:   # the current version of the project
  language:  # the lanauge (ruby, golang, etc.) of the engine
  summary:   # a 140 character summary of the project

Please ensure all required fields are provided and try again.`))

		os.Exit(1)
	}

	// attempt to read a README.md file and add it to the release...
	b, err := ioutil.ReadFile("./README.md")
	if err != nil {

		// this only fails if the file is not found, EOF is not an error. If no Readme
		// is found exit 1
		fmt.Printf(stylish.Error("missing readme", "Your engine is missing a README.md file. This file is required for publishing, as it is the only way for you to communicate how to use your engine. Please add a README.md and try again."))
		os.Exit(1)
	}

	//
	release.Readme = string(b)

	// check to see if the engine already exists on nanobox.io
	fmt.Printf(stylish.Bullet("Checking for existing engine on nanobox.io"))
	engine, err := api.GetEngine(api.UserSlug, release.Name)

	// if no engine exists, create a new one
	if err != nil {

		// if no engine is found create one
		if apiErr, _ := err.(api.APIError); apiErr.Code == 404 {

			fmt.Printf(stylish.SubTaskStart("Creating new engine on nanobox.io"))

			//
			engine = &api.Engine{
				Generic:      release.Generic,
				LanguageName: release.Language,
				Name:         release.Name,
			}

			//
			if _, err := api.CreateEngine(engine); err != nil {
				fmt.Printf(stylish.ErrBullet("Unable to create engine (%v).", err))
				os.Exit(1)
			}

			// wait until engine has been successfuly created before uploading to s3
			for {
				fmt.Print(".")

				p, err := api.GetEngine(api.UserSlug, release.Name)
				if err != nil {
					config.Fatal("[commands/publish] api.GetEngine failed", err.Error())
				}

				// once the engine is "active", break
				if p.State == "active" {
					break
				}

				//
				time.Sleep(1000 * time.Millisecond)
			}

			// generically handle any other errors
		} else {
			config.Fatal("[commands/publish] api.GetEngine failed", err.Error())
		}

		stylish.Success()
	}

	// create a meta.json file where we can add any extra data we might need
	meta, err := os.Create("./meta.json")
	if err != nil {
		config.Fatal("[commands/publish] os.Create() failed", err.Error())
	}
	defer meta.Close()
	defer os.Remove(meta.Name())

	//
	meta.WriteString(fmt.Sprintf(`{"engine_id": "%s"}`, engine.ID))

	// this is our predefined list of everything that gets archived as part of the
	// engine being published
	files := map[string][]string{
		"required": []string{"./bin", "./Enginefile", "./meta.json"},
		"optional": []string{"./lib", "./templates", "./files"},
	}

	//
	for k, v := range files {
		if k == "required" {

			// check to ensure no required files are missing
			for _, f := range v {
				if _, err := os.Stat(f); err != nil {
					fmt.Printf(stylish.Error("required files missing", "Your Engine is missing one or more required files for publishing. Please read the following documentation to ensure all required files are included and try again.:\n\ndocs.nanobox.io/engines/project-creation/#example-engine-file-structure\n"))
					os.Exit(1)
				}
			}
		}
	}

	// once the whole thing is working again, try swaping the go routine to be on
	// readers instead of the writer. the writer will block until readers are done
	// reading, so there may not be a need for the wait groups.

	// write the archive to a local file
	// archive, err := os.Create(fmt.Sprintf("%v-%v.release.tgz", release.Name, release.Version))
	// if err != nil {
	// 	config.Fatal("[commands/publish] os.Create() failed", err.Error())
	// }
	// defer archive.Close()

	// create an empty buffer for writing the file contents to for the subsequent
	// upload
	archive := bytes.NewBuffer(nil)

	//
	h := md5.New()

	//
	mw := io.MultiWriter(h, archive)

	//
	gzw := gzip.NewWriter(mw)

	//
	tw = tar.NewWriter(gzw)

	//
	wg := &sync.WaitGroup{}
	wg.Add(1)

	//
	go func() {

		defer gzw.Close()
		defer tw.Close()

		// range over each file type...
		for _, v := range files {

			// range over each file of each type...
			for _, f := range v {

				// required files have alrady been checked, so skip any remaining (optional)
				// files/folders that arent here
				if _, err := os.Stat(f); err != nil {
					continue
				}

				// tarball any remaining files/folders that are found
				if err := filepath.Walk(f, tarFile); err != nil {
					config.Fatal("[commands/publish] filepath.Walk() failed", err.Error())
				}
			}
		}

		wg.Done()
	}()

	wg.Wait()

	// add the checksum for the new release once its finished being archived
	release.Checksum = fmt.Sprintf("%x", h.Sum(nil))

	//
	// attempt to upload the release to S3

	//
	fmt.Printf(stylish.Bullet("Uploading release to s3..."))

	v := url.Values{}
	v.Add("user_slug", api.UserSlug)
	v.Add("auth_token", api.AuthToken)
	v.Add("version", release.Version)

	//
	s3url, err := s3.RequestURL(fmt.Sprintf("http://api.nanobox.io/v1/engines/%v/request_upload?%v", release.Name, v.Encode()))
	if err != nil {
		config.Fatal("[commands/publish] s3.RequestURL() failed", err.Error())
	}

	//
	if err := s3.Upload(s3url, archive); err != nil {
		config.Fatal("[commands/publish] s3.Upload() failed", err.Error())
	}

	//
	// if the release uploaded successfully to s3, created one on odin
	fmt.Printf(stylish.Bullet("Uploading release to nanobox.io"))
	if _, err := api.CreateEngineRelease(release.Name, release); err != nil {
		fmt.Printf(stylish.ErrBullet("Unable to publish release (%v).", err))
		os.Exit(1)
	}
}

// tarFile
func tarFile(path string, fi os.FileInfo, err error) error {

	// only want to tar files...
	if !fi.Mode().IsDir() {

		// fmt.Println("TARING!", path)

		// create header for this file
		header := &tar.Header{
			Name:    path,
			Size:    fi.Size(),
			Mode:    int64(fi.Mode()),
			ModTime: fi.ModTime(),
		}

		// write the header to the tarball archive
		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		// open the file for taring...
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		// copy the file data to the tarball
		if _, err := io.Copy(tw, f); err != nil {
			return err
		}
	}

	return nil
}
