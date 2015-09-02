// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

package commands

//
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

	api "github.com/pagodabox/nanobox-api-client"
	"github.com/pagodabox/nanobox-cli/auth"
	"github.com/pagodabox/nanobox-cli/config"
	"github.com/pagodabox/nanobox-cli/util"
	"github.com/pagodabox/nanobox-golang-stylish"
)

var tw *tar.Writer

//
var publishCmd = &cobra.Command{
	Use:   "publish",
	Short: "Publishes an engine to nanobox.io",
	Long: `
Description:
  Publishes an engine to nanobox.io`,

	Run: nanoPublish,
}

// nanoPublish
func nanoPublish(ccmd *cobra.Command, args []string) {

	//
	api.UserSlug, api.AuthToken = auth.Authenticate()

	//
	stylish.Header("publishing engine")

	// attempt to parse an enginefile
	if err := config.Enginefile.Parse(); err != nil {
		util.LogFatal("commands/init] config.Enginefile.Parse() failed", err)
	}

	// create a new release based off the enginefile config options
	fmt.Printf(stylish.Bullet("Creating release..."))
	release := &api.EngineReleaseCreateOptions{
		Authors:      config.Enginefile.Authors,
		Description:  config.Enginefile.Description,
		License:      config.Enginefile.License,
		Name:         config.Enginefile.Name,
		ProjectFiles: config.Enginefile.ProjectFiles,
		Readme:       config.Enginefile.Readme,
		Stability:    config.Enginefile.Stability,
		Summary:      config.Enginefile.Summary,
		Version:      config.Enginefile.Version,
	}

	// GET to API to see if engine exists
	fmt.Printf(stylish.Bullet("Checking for existing engine on nanobox.io"))
	if _, err := api.GetEngine(api.UserSlug, release.Name); err != nil {

		// if no engine is found create one
		if apiErr, _ := err.(api.APIError); apiErr.Code == 404 {
			stylish.SubTaskStart("No engine found, creating new engine on nanobox.io...")

			//
			engineCreateOptions := &api.EngineCreateOptions{Name: release.Name}
			if _, err := api.CreateEngine(engineCreateOptions); err != nil {
				util.LogFatal("[commands.publish] api.CreateEngine() failed", err)
			}

			// wait until engine has been successfuly created before uploading to s3
			for {
				fmt.Print(".")

				p, err := api.GetEngine(api.UserSlug, release.Name)
				if err != nil {
					util.LogFatal("[commands.publish] api.GetEngine() failed", err)
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
			util.LogFatal("[commands publish] api.GetEngine failed", err)
		}

		stylish.Success()
	}

	// once the whole thing is working again, try swaping the go routine to be on
	// readers instead of the writer. the writer will block until readers are done
	// reading, so there may not be a need for the wait groups.

	// archive, err := os.Create(fmt.Sprintf("%v-%v.release.tgz", release.Name, release.Version))
	// if err != nil {
	// 	util.LogFatal("[commands.publish] os.Create() failed", err)
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

		for _, files := range release.ProjectFiles {
			if err := filepath.Walk(files, tarFile); err != nil {
				util.LogFatal("[commands.publish] filepath.Walk() failed", err)
			}
		}

		wg.Done()
	}()

	wg.Wait()

	//
	fmt.Printf(stylish.Bullet("Uploading release to s3..."))

	v := url.Values{}
	v.Add("user_slug", api.UserSlug)
	v.Add("auth_token", api.AuthToken)
	v.Add("version", release.Version)

	//
	s3url, err := util.RequestS3URL(fmt.Sprintf("http://api.nanobox.io/v1/engines/%v/request_upload?%v", release.Name, v.Encode()))
	if err != nil {
		util.LogFatal("[commands/publish] util.RequestS3URL failed", err)
	}

	//
	if err := util.S3Upload(s3url, archive); err != nil {
		util.LogFatal("[commands/publish] util.S3Upload failed", err)
	}

	// add readme to release (if any)
	b, err := ioutil.ReadFile(release.Readme)
	if err != nil {
		config.Console.Warn("No readme found at '%v', continuing...", release.Readme)
	}

	// prepare the release for upload
	release.Checksum = fmt.Sprintf("%x", h.Sum(nil))
	release.ProjectFiles = append(release.ProjectFiles, "Enginefile")
	release.Readme = string(b)

	// if the release uploaded successfully to s3, created one on odin
	fmt.Printf(stylish.Bullet("Uploading release to nanobox.io"))
	if _, err := api.CreateEngineRelease(release.Name, release); err != nil {
		util.LogFatal("[commands.publish] api.CreateEngineRelease() failed", err)
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
