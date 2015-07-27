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
	"github.com/pagodabox/nanobox-golang-stylish"
)

type (

	// PublishCommand satisfies the Command interface for listing a user's apps
	PublishCommand struct {
		tw *tar.Writer
	}

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

// Help
func (c *PublishCommand) Help() {
	ui.CPrint(`
Description:
  Publish your engine to nanobox.io

Usage:
  nanobox publish
  `)
}

// Run
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

	stylish.Header("publishing engine")

	// look for an Enginefile to parse
	fmt.Printf(stylish.Bullet("Parsing Enginefile"))
	fi, err := os.Stat("./Enginefile")
	if err != nil {
		fmt.Println("Enginefile not found. Be sure to publish from a project directory. Exiting... ")
		config.Log.Fatal("[commands.publish] os.Stat() failed", err)
		os.Exit(1)
	}

	//
	file, err := ioutil.ReadFile(fi.Name())
	if err != nil {
		ui.LogFatal("[commands.publish] ioutil.ReadFile() failed", err)
	}

	//
	fmt.Printf(stylish.Bullet("Creating release"))
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

	// add Enginefile to the list of a releases project_files
	release.ProjectFiles = append(release.ProjectFiles, "Enginefile")

	// GET to API to see if engine exists
	fmt.Printf(stylish.Bullet("Checking for existing engine on nanobox.io"))
	engine, err := api.GetEngine(api.UserSlug, release.Name)
	if err != nil {
		fmt.Printf(stylish.Bullet("No engine found..."))
	}

	// if no engine is found create one
	if engine.ID == "" {
		stylish.SubTaskStart("Creating new engine on nanobox.io")

		engineCreateOptions := &api.EngineCreateOptions{Name: release.Name, Type: release.Type}
		if _, err := api.CreateEngine(engineCreateOptions); err != nil {
			ui.LogFatal("[commands.publish] api.CreateEngine() failed", err)
		}

		// wait until engine has been successfuly created before uploading to warehouse
		for {
			fmt.Print(".")

			p, err := api.GetEngine(api.UserSlug, release.Name)
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

		stylish.Success()
	} else {
		config.Console.Info("Engine found on nanobox by the name '%v'", release.Name)
	}

	//
	h := md5.New()

	//
	pr, pw := io.Pipe()

	//
	mw := io.MultiWriter(h, pw)

	//
	gzw := gzip.NewWriter(mw)

	//
	c.tw = tar.NewWriter(gzw)

	//
	wg := &sync.WaitGroup{}

	wg.Add(1)

	//
	go func() {

		// defer is LIFO
		defer pw.Close()
		defer gzw.Close()
		defer c.tw.Close()

		for _, pf := range release.ProjectFiles {

			if stat, err := os.Stat(pf); err == nil {

				// if its a directory, walk the directory taring each file
				if stat.Mode().IsDir() {
					if err := filepath.Walk(pf, c.tarDir); err != nil {
						ui.LogFatal("[commands.publish] filepath.Walk() failed", err)
					}

					// if its a file tar it
				} else {
					c.tarFile(pf)
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
		"Objectalias": "releases/" + release.Version,
	}

	// upload tarball release to warehouse
	fmt.Printf(stylish.Bullet("Uploading release to nanobox warehouse..."))

	if err := api.DoRawRequest(obj, "POST", "http://warehouse.nanobox.io/objects", pr, headers); err != nil {
		ui.LogFatal("[commands.publish] api.DoRawRequest() failed", err)
	}

	wg.Wait()

	defer pr.Close()

	checksum := fmt.Sprintf("%x", h.Sum(nil))

	// check checksum
	if checksum == obj.CheckSum {

		release.Checksum = checksum

		fmt.Printf(stylish.Bullet("Uploading release to nanobox.io"))

		// POST release on API (odin)
		if _, err := api.CreateEngineRelease(engine.Name, release); err != nil {
			ui.LogFatal("[commands.publish] api.CreateEngineRelease() failed", err)
		}
	} else {
		config.Console.Fatal("Checksums don't match!!! Exiting...")
		os.Exit(1)
	}
}

// tarDir
func (c *PublishCommand) tarDir(path string, fi os.FileInfo, err error) error {
	if fi.Mode().IsDir() {
		return nil
	}

	if err := c.tarFile(path); err != nil {
		return err
	}

	return nil
}

// tarFile
func (c *PublishCommand) tarFile(path string) error {

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
		if err := c.tw.WriteHeader(header); err != nil {
			return err
		}

		// copy the file data to the tarball
		if _, err := io.Copy(c.tw, f); err != nil {
			return err
		}
	}

	return nil
}
