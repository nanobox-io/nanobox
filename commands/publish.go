// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

package commands

import (
	"archive/tar"
	"compress/gzip"
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/ghodss/yaml"

	api "github.com/pagodabox/nanobox-api-client"
	"github.com/pagodabox/nanobox-cli/auth"
	"github.com/pagodabox/nanobox-cli/config"
	"github.com/pagodabox/nanobox-cli/ui"
)

type (

	// PublishCommand satisfies the Command interface for listing a user's apps
	PublishCommand struct{}

	// Object is returned from warehouse
	Object struct {
		Alias    string
		BucketID string
		CheckSum string
		ID       string
		Public   bool
		Size     int64
	}
)

// Help prints detailed help text for the app list command
func (c *PublishCommand) Help() {
	ui.CPrint(`
Description:
  Publish your engine to nanobox.io

Usage:
  nanobox publish
  `)
}

//
var tw *tar.Writer

// Run displays select information about all of a user's apps
func (c *PublishCommand) Run(opts []string) {

	// check for auth
	if !auth.IsAuthenticated() {
		fmt.Println("Before using the Pagoda Box CLI on this machine, please login to your account:")

		userslug := ui.Prompt("Username: ")
		password := ui.PPrompt("Password: ")

		// authenticate
		if err := auth.Authenticate(userslug, password); err != nil {
			ui.LogFatal("[main] auth.Authenticate() failed", err)
		}

		fmt.Println("To begin using the Pagoda Box CLI type 'pagoda' to see a list of commands.")
		os.Exit(0)
	}

	// pull the users api credentials
	creds, err := auth.Credentials()
	if err != nil {
		ui.LogFatal("[main] auth.Credentials() failed", err)
	}

	api.UserSlug = creds["user_slug"]
	api.AuthToken = creds["auth_token"]

	// if the credentials are empty attempt to reauthenticate
	if api.UserSlug == "" || api.AuthToken == "" {
		config.Console.Warn("No login credentials found! Reauthenticating...")
		auth.ReAuthenticate()
	}

	// look for a Enginefile to parse
	enginefile, err := os.Stat("./Enginefile")
	if err != nil {
		fmt.Println("Enginefile not found. Be sure to publish from a project directory. Exiting... ")
		config.Log.Fatal("[commands.publish] os.Stat() failed", err)
		os.Exit(1)
	}

	//
	file, err := ioutil.ReadFile(enginefile.Name())
	if err != nil {
		ui.LogFatal("[commands.publish] ioutil.ReadFile() failed", err)
	}

	//
	release := &api.EngineReleaseCreateOptions{}
	if err := yaml.Unmarshal(file, release); err != nil {
		ui.LogFatal("[commands.publish] yaml.Unmarshal() failed", err)
	}

	// add readme to release
	b, err := ioutil.ReadFile(release.Readme)
	if err != nil {
		config.Console.Warn("No readme found at '%v', continuing...", release.Readme)
	}

	release.Readme = string(b)

	// GET to API to see if engine exists
	engine, err := api.GetEngine(release.Name)
	if err != nil {
		fmt.Println("ERR!!!", err)
		config.Console.Info("No engines found on nanobox by the name '%v'", release.Name)
	}

	// if no engine is found create one
	if engine.ID == "" {
		engineCreateOptions := &api.EngineCreateOptions{Name: release.Name, Type: release.Type}
		if _, err := api.CreateEngine(engineCreateOptions); err != nil {
			ui.LogFatal("[commands.publish] api.CreateEngine() failed", err)
		}

		fmt.Print("Creating engine..")

		// wait until engine has been successfuly created before uploading to warehouse
		for {
			fmt.Print(".")

			p, err := api.GetEngine(release.Name)
			if err != nil {
				ui.LogFatal("[commands.publish] api.GetEngine() failed", err)
			}

			// once the service has a tunnel ip and port break
			if p.State == "active" {

				// set our engine to the active one
				engine = p
				break
			}

			//
			time.Sleep(1000 * time.Millisecond)
		}

		fmt.Println(" complete")
	} else {
		config.Console.Info("Engine found on nanobox by the name '%v'", release.Name)
	}

	// upload tarball release to warehouse
	config.Console.Info("Uploading release to warehouse...")

	//
	h := md5.New()

	//
	pr, pw := io.Pipe()

	//
	mw := io.MultiWriter(h, pw)

	//
	gzw := gzip.NewWriter(mw)

	//
	tw = tar.NewWriter(gzw)

	//
	wg := &sync.WaitGroup{}

	wg.Add(1)

	//
	go func() {

		// defer is LIFO
		defer pw.Close()
		defer gzw.Close()
		defer tw.Close()

		for _, pf := range release.ProjectFiles {

			if stat, err := os.Stat(pf); err == nil {

				// if its a directory, walk the directory taring each file
				if stat.Mode().IsDir() {
					if err := filepath.Walk(pf, tarDir); err != nil {
						ui.LogFatal("[commands.publish] filepath.Walk() failed", err)
					}

					// if its a file tar it
				} else {
					tarFile(pf)
				}
			}
		}

		wg.Done()
	}()

	obj := &Object{}

	//
	headers := map[string]string{
		"Userid":      engine.WarehouseUser,
		"Key":         engine.WarehouseKey,
		"Bucketid":    engine.ID,
		"Objectalias": "release-" + release.Version,
	}

	//
	if err := api.DoRawRequest(obj, "POST", "http://warehouse.nanobox.io/objects", pr, headers); err != nil {
		ui.LogFatal("[commands.publish] api.DoRawRequest() failed", err)
	}

	wg.Wait()

	defer pr.Close()

	checksum := fmt.Sprintf("%x", h.Sum(nil))

	// check checksum
	if checksum == obj.CheckSum {

		release.Checksum = checksum

		// POST release on API (odin)
		if _, err := api.CreateEngineRelease(engine.Name, release); err != nil {
			ui.LogFatal("[commands.publish] api.CreateEngineRelease() failed", err)
		}
	} else {
		config.Console.Fatal("Checksums don't match!!! Exiting...")
		os.Exit(1)
	}
}

//
func tarDir(path string, fi os.FileInfo, err error) error {
	if fi.Mode().IsDir() {
		return nil
	}

	if err := tarFile(path); err != nil {
		return err
	}

	return nil
}

//
func tarFile(path string) error {

	// open the file/dir...
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	// stat the file
	if fi, err := f.Stat(); err == nil {

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

		// copy the file data to the tarball
		if _, err := io.Copy(tw, f); err != nil {
			return err
		}
	}

	return nil
}
