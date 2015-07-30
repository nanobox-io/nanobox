// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

package config

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	semver "code.google.com/p/go-semver/version"

	"github.com/jcelliott/lumber"
	"github.com/mitchellh/go-homedir"
)

//
const (
	VERSION = "0.0.9"
)

//
var (
	App        string
	AppDir     string
	AppsDir    string
	AuthFile   string
	CWDir      string
	LogFile    string
	HomeDir    string
	NanoDir    string
	UpdateFile string

	//
	Console  *lumber.ConsoleLogger
	Log      *lumber.FileLogger
	LogLevel int

	//
	Version *semver.Version

	//
	Boxfile  *BoxfileConfig
	Nanofile *NanofileConfig
)

// Init sets up a HomeDir, and NanoDir
func init() {

	// set the default log level
	LogLevel = lumber.INFO

	// check for debug mode and set the appropriate log level
	if os.Args[len(os.Args)-1] == "--debug" {
		LogLevel = lumber.DEBUG
	}

	//
	Console = lumber.NewConsoleLogger(LogLevel)

	//
	homeDir, err := homedir.Dir()
	if err != nil {
		fmt.Println("Fatal error! See ~/.nanobox/nanobox.log for details. Exiting...")
		Log.Fatal("[config] homedir.Dir() failed %v\n", err)
		Log.Close()
		os.Exit(1)
	}

	HomeDir = homeDir
	NanoDir = path.Clean(HomeDir + "/.nanobox")
	AuthFile = filepath.Clean(NanoDir + "/.auth")
	LogFile = path.Clean(NanoDir + "/nanobox.log")
	UpdateFile = path.Clean(NanoDir + "/.update")

	// get the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("BONK!", err)
	}

	CWDir = cwd

	// the 'app' name is the base folder of the cwd
	App = path.Base(cwd)
	AppsDir = path.Clean(NanoDir + "/apps")
	AppDir = fmt.Sprintf("%s/%s", AppsDir, App)

	//
	version, err := semver.Parse(VERSION)
	if err != nil {
		fmt.Println("Fatal error! See ~/.nanobox/nanobox.log for details. Exiting...")
		Log.Fatal("[config] semver.Parse() failed %v", err)
		Log.Close()
		os.Exit(1)
	}

	Version = version

}
