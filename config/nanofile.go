// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

//
package config

import (
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
)

// NanofileConfig represents all available/expected .nanofile configurable options
type NanofileConfig struct {
	CPUCap   int    `json:"cpu_cap"`  // max %CPU usage allowed to the guest vm
	CPUs     int    `json:"cpus"`     // number of CPUs to dedicate to the guest vm
	Domain   string `json:"domain"`   // the domain to use in conjuntion with the ip when accesing the guest vm (defaults to <Name>.dev)
	IP       string `json:"ip"`       // the ip added to the /etc/hosts file for accessing the guest vm
	Name     string `json:"name"`     // the name given to the project (defaults to cwd)
	Provider string `json:"provider"` // guest vm provider (virtual box, vmware, etc)
	RAM      int    `json:"ram"`      // ammount of RAM to dedicate to the guest vm
}

// ParseNanofile
func ParseNanofile() *NanofileConfig {

	//
	nanofile := &NanofileConfig{
		CPUCap:   50,
		CPUs:     2,
		Name:     filepath.Base(CWDir),
		Provider: "virtualbox",
		RAM:      1024,
	}

	nanofilePath := Root + "/.nanofile"

	// look for a global .nanofile first in the ~/.nanobox directory, and override
	// any default options found.
	if _, err := os.Stat(nanofilePath); err == nil {
		if err := ParseConfig(nanofilePath, nanofile); err != nil {
			fmt.Printf("Nanobox failed to parse your .nanofile. Please ensure it is valid YAML and try again.\n")
			Exit(1)
		}
	}

	nanofilePath = "./.nanofile"

	// then look for a local .nanofile and override any global, or remaining default
	// options found
	if _, err := os.Stat(nanofilePath); err == nil {
		if err := ParseConfig(nanofilePath, nanofile); err != nil {
			fmt.Printf("Nanobox failed to parse your .nanofile. Please ensure it is valid YAML and try again.\n")
			Exit(1)
		}
	}

	// set name specific options after potential .nanofiles have been parsed
	nanofile.Domain = fmt.Sprintf("%s.dev", nanofile.Name)

	// assign a default IP if none is specified
	if nanofile.IP == "" {
		nanofile.IP = appNameToIP(nanofile.Name)
	}

	return nanofile
}

// appNameToIP generates an IPv4 address based off the app name for use as a
// vagrant private_network IP.
func appNameToIP(name string) string {

	var network uint32 = 2886729728 // 172.16.0.0 network
	var sum uint32 = 0              // the last two octets of the assigned network

	// create an md5 of the app name to ensure a uniqe IP is generated each time
	h := md5.New()
	io.WriteString(h, name)

	// iterate through each byte in the md5 hash summing along the way
	for _, v := range []byte(h.Sum(nil)) {
		sum += uint32(v)
	}

	ip := make(net.IP, 4)

	// convert app name into a private network IP by adding the first portion of
	// the network with the generated portion
	binary.BigEndian.PutUint32(ip, (network + sum))

	return ip.String()
}
