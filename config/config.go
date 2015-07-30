// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

package config

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"path"
	"path/filepath"

	semver "code.google.com/p/go-semver/version"

	"github.com/jcelliott/lumber"
	"github.com/mitchellh/go-homedir"
)

//
const (
	VERSION = "0.1.4"
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

// Parser
type Parser interface {
	Parse() error //
}

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
		Log.Fatal("[config] semver.Parse() failed", err)
		Log.Close()
		os.Exit(1)
	}

	Version = version

	// create a default BoxfileConfig
	Boxfile = &BoxfileConfig{}

	// create a default NanofileConfig
	Nanofile = &NanofileConfig{
		CPUCap:   50,
		CPUs:     2,
		Domain:   fmt.Sprintf("%v.nano.dev", App),
		IP:       appNameToIP(App),
		Provider: "virtualbox",
		RAM:      1024,
	}
}

// appNameToIP generates an IPv4 address based off the app name for use as a
// vagrant private_network IP.
func appNameToIP(name string) string {

	var sum uint32 = 0
	var network uint32 = 2886729728 // 172.16.0.0 network

	for _, value := range []byte(name) {
		sum += uint32(value)
	}

	ip := make(net.IP, 4)

	// convert app name into a private network IP
	binary.BigEndian.PutUint32(ip, (network + sum))

	return ip.String()
}
