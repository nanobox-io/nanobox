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
)

// Help prints detailed help text for the app list command
func (c *PublishCommand) Help() {
	ui.CPrint(`
Description:
  Publish your package to nanobox.io

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

	// look for a Packagefile to parse
	pf, err := os.Stat("./Packagefile")
	if err != nil {
		fmt.Println("Packagefile not found. Be sure to publish from a project directory. Exiting... ")
		config.Log.Fatal("[commands.publish] os.Stat() failed %v", err)
		os.Exit(1)
	}

	//
	file, err := ioutil.ReadFile(pf.Name())
	if err != nil {
		ui.LogFatal("[commands.publish] ioutil.ReadFile() failed: %v", err)
	}

	//
	release := &api.Release{}
	if err := yaml.Unmarshal(file, release); err != nil {
		ui.LogFatal("[commands.publish] yaml.Unmarshal() failed: %v", err)
	}

	// add readme to release
	b, err := ioutil.ReadFile(release.Readme)
	if err != nil {
		config.Console.Warn("No readme found at '%v', continuing...", release.Readme)
	}

	release.Readme = string(b)

	// set up the output file

	archive := "release.tgz"

	a, err := os.Create(archive)
	if err != nil {
		ui.LogFatal("[commands.publish] os.Create() failed: %v", err)
	}
	defer a.Close()

	// set up the gzip writer
	gw := gzip.NewWriter(a)
	defer gw.Close()

	// set up the tar writer
	tw = tar.NewWriter(gw)
	defer tw.Close()

	//
	for _, dep := range release.Dependencies {

		if stat, err := os.Stat(dep); err == nil {

			// if its a directory, walk the directory taring each file
			if stat.Mode().IsDir() {
				if err := filepath.Walk(dep, tarDir); err != nil {
					ui.LogFatal("[commands.publish] filepath.Walk() failed: %v", err)
				}

				// if its a file tar it
			} else {
				tarFile(dep)
			}
		}
	}

	// create a checksum of tarball
	h := md5.New()
	tarball, _ := ioutil.ReadFile(archive)
	h.Write(tarball)

	// GET to API to see if project exists
	pkg, err := api.GetPackage(release.Name)
	if err != nil {
		config.Console.Info("No packages found on nanobox by the name '%v'", release.Name)
	}

	// if no package is found create one
	if pkg.ID == "" {
		packageCreateOptions := &api.PackageCreateOptions{Name: release.Name, Type: release.Type}
		if _, err := api.CreatePackage(packageCreateOptions); err != nil {
			ui.LogFatal("[commands.publish] api.CreatePackage() failed: %v", err)
		}

		fmt.Print("Creating package..")

		// wait until package has been successfuly created before uploading to warehouse
		for {
			fmt.Print(".")

			p, err := api.GetPackage(release.Name)
			if err != nil {
				ui.LogFatal("[commands.publish] api.GetPackage() failed: %v", err)
			}

			// once the service has a tunnel ip and port break
			if p.State == "active" {
				break
			}

			//
			time.Sleep(1000 * time.Millisecond)
		}

		fmt.Println(" complete")
	}

	// upload tarball release to warehouse
	config.Console.Info("Uploading release to warehouse...")

	// POST release on API (odin)
	releaseCreateOptions := &api.ReleaseCreateOptions{}
	if _, err := api.CreateRelease(releaseCreateOptions); err != nil {
		ui.LogFatal("[commands.publish] api.CreateRelease() failed: ", err)
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
