// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

//
package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"github.com/ghodss/yaml"
	"github.com/mitchellh/go-homedir"

	"github.com/nanobox-io/nanobox-golang-stylish"
)

//
const (
	VERSION = "0.15.7"

	SERVER_PORT = ":1757"
	MIST_PORT   = ":1445"
	LOGTAP_PORT = ":6361"
)

type (
	exiter func(int)
)

//
var (
	err   error //
	mutex *sync.Mutex

	//
	AppDir     string // the path to the application (~.nanobox/apps/<app>)
	AppsDir    string // ~.nanobox/apps
	CWDir      string // the current working directory
	EnginesDir string // ~.nanobox/engines
	Home       string // the users home directory (~)
	IP         string // the guest vm's private network ip (generated from app name)
	Root       string // nanobox's root directory path (~.nanobox)
	UpdateFile string // the path to the .update file (~.nanobox/.update)

	//
	Nanofile *NanofileConfig // parsed nanofile options
	VMfile   *VMfileConfig   // parsed nanofile options

	//
	ServerURI string // nanobox-server host:port combo (IP:1757)
	ServerURL string // nanobox-server host:port combo (IP:1757) (http)
	MistURI   string // mist's host:port combo (IP:1445)
	LogtapURI string // logtap's host:port combo (IP:6361)

	// flags
	Background bool   //
	Devmode    bool   //
	Force      bool   //
	Verbose    bool   //
	Silent     bool   //
	LogLevel   string //

	//
	Exit exiter = os.Exit
)

//
func init() {

	// default log level
	LogLevel = "info"

	// set the current working directory first, as it's used in other steps of the
	// configuration process
	if CWDir, err = os.Getwd(); err != nil {
		Log.Fatal("[config/config] os.Getwd() failed", err.Error())
	}

	// set Home based off the users homedir (~)
	if Home, err = homedir.Dir(); err != nil {
		Log.Fatal("[config/config] homedir.Dir() failed", err.Error())
	}

	// set nanobox's root directory;
	Root = filepath.Join(Home, ".nanobox")

	// check for a ~/.nanobox dir and create one if it's not found
	if _, err := os.Stat(Root); err != nil {
		fmt.Printf(stylish.Bullet("Creating %s directory", Root))
		if err := os.Mkdir(Root, 0755); err != nil {
			Log.Fatal("[config/config] os.Mkdir() failed", err.Error())
		}
	}

	// check for a ~/.nanobox/.update file and create one if it's not found
	UpdateFile = filepath.Join(Root, ".update")
	if _, err := os.Stat(UpdateFile); err != nil {
		f, err := os.Create(UpdateFile)
		if err != nil {
			Log.Fatal("[config/config] os.Create() failed - ", err.Error())
		}
		defer f.Close()
	}

	// check for a ~/.nanobox/engines dir and create one if it's not found
	EnginesDir = filepath.Join(Root, "engines")
	if _, err := os.Stat(EnginesDir); err != nil {
		if err := os.Mkdir(EnginesDir, 0755); err != nil {
			Log.Fatal("[config/config] os.Mkdir() failed", err.Error())
		}
	}

	// check for a ~/.nanobox/apps dir and create one if it's not found
	AppsDir = filepath.Join(Root, "apps")
	if _, err := os.Stat(AppsDir); err != nil {
		if err := os.Mkdir(AppsDir, 0755); err != nil {
			Log.Fatal("[config/config] os.Mkdir() failed", err.Error())
		}
	}

	// the .nanofile needs to be parsed right away so that its config options are
	// available as soon as possible
	Nanofile = ParseNanofile()

	//
	ServerURI = Nanofile.IP + SERVER_PORT
	ServerURL = "http://" + ServerURI
	MistURI = Nanofile.IP + MIST_PORT
	LogtapURI = Nanofile.IP + LOGTAP_PORT

	// set the 'App' first so it can be used in subsequent configurations; the 'App'
	// is set to the name of the cwd; this can be overriden from a .nanofile
	AppDir = filepath.Join(AppsDir, Nanofile.Name)
}

// ParseConfig
func ParseConfig(path string, v interface{}) error {

	//
	fp, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	//
	f, err := ioutil.ReadFile(fp)
	if err != nil {
		return err
	}

	//
	return yaml.Unmarshal(f, v)
}

// writeConfig
func writeConfig(path string, v interface{}) error {

	// take a config objects path and create (and truncate) the file, preparing it
	// to receive new configurations
	f, err := os.Create(path)
	if err != nil {
		Fatal("[config/config] os.Create() failed", err.Error())
	}
	defer f.Close()

	// marshal the config object
	b, err := yaml.Marshal(v)
	if err != nil {
		Fatal("[config/config] yaml.Marshal() failed", err.Error())
	}

	// mutex.Lock()

	// write it back to the file
	if _, err := f.Write(b); err != nil {
		return err
	}

	// mutex.Unlock()

	return nil
}
