// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

package config

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"sync"

	"github.com/ghodss/yaml"
	"github.com/mitchellh/go-homedir"

	"github.com/pagodabox/nanobox-golang-stylish"
)

//
const (
	SERVER_PORT = ":1757"
	MIST_PORT   = ":1445"
	LOGTAP_PORT = ":6361"
)

//
var (
	err   error //
	mutex *sync.Mutex

	Lock net.Conn //

	//
	AppDir string // the path to the application (~.nanobox/apps/<app>)
	CWDir  string // the current working directory
	Home   string // the users home directory (~)
	IP     string // the guest vm's private network ip (generated from app name)
	Root   string // nanobox's root directory path (~.nanobox)

	//
	Nanofile *NanofileConfig // parsed nanofile options
	VMfile   *VMfileConfig   // parsed nanofile options

	//
	ServerURI string // nanobox-server host:port combo (IP:1757)
	ServerURL string // nanobox-server host:port combo (IP:1757) (http)
	MistURI   string // mist's host:port combo (IP:1445)
	LogtapURI string // logtap's host:port combo (IP:6361)
)

//
type (

	// BoxfileConfig represents all available/expected Boxfile configurable options
	BoxfileConfig struct {
		Build Build //
	}

	// Build represents a possible node in the Boxfile with it's own set of options
	Build struct {
		Engine string `json:"engine"` //
	}

	// NanofileConfig represents all available/expected .nanofile configurable options
	NanofileConfig struct {
		CPUCap   int    `json:"cpu_cap"`  // max %CPU usage allowed to the guest vm
		CPUs     int    `json:"cpus"`     // number of CPUs to dedicate to the guest vm
		Domain   string `json:"domain"`   // the domain to use in conjuntion with the ip when accesing the guest vm (defaults to <Name>.dev)
		IP       string `json:"ip"`       // the ip added to the /etc/hosts file for accessing the guest vm
		Name     string `json:"name"`     // the name given to the project (defaults to cwd)
		Provider string `json:"provider"` // guest vm provider (virtual box, vmware, etc)
		RAM      int    `json:"ram"`      // ammount of RAM to dedicate to the guest vm
	}

	// VMfileConfig represents all available/expected .vmfile configurable options
	VMfileConfig struct {
		Deployed    bool   // was the most recent deploy successufl
		Mode        string // foreground/background
		Status      string // the current staus of the VM
		Suspendable bool   // is the VM able to be suspended
		UUID        string // the UUID of the VM
	}
)

//
func init() {

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
	Root = filepath.Clean(Home + "/.nanobox")

	// check for a ~/.nanobox dir and create one if it's not found
	if _, err := os.Stat(Root); err != nil {
		fmt.Printf(stylish.Bullet("Creating %s directory", Root))
		if err := os.Mkdir(Root, 0755); err != nil {
			Log.Fatal("[config/config] os.Mkdir() failed", err.Error())
		}
	}

	// check for a ~/.nanobox/apps dir and create one if it's not found
	apps := filepath.Clean(Root + "/apps")
	if _, err := os.Stat(apps); err != nil {
		if err := os.Mkdir(apps, 0755); err != nil {
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
	AppDir = apps + "/" + Nanofile.Name
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
	if err := yaml.Unmarshal(f, v); err != nil {
		return err
	}

	return nil
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
		return err
	}

	// mutex.Lock()

	// write it back to the file
	if _, err := f.Write(b); err != nil {
		return err
	}

	// mutex.Unlock()

	return nil
}
