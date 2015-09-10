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
	"path/filepath"
)

// ParseNanofile
func ParseNanofile() *NanofileConfig {

	//
	nanofile := &NanofileConfig{
		path:     Root + "/.nanofile",
		CPUCap:   50,
		CPUs:     2,
		Name:     filepath.Base(CWDir),
		Provider: "virtualbox",
		RAM:      1024,
	}

	// look for a global .nanofile first in the ~/.nanobox directory, and override
	// any default options found.
	if fi, _ := os.Stat(nanofile.path); fi != nil {
		if err := ParseConfig(nanofile.path, nanofile); err != nil {
			fmt.Printf("Nanobox failed to parse your .nanofile. Please ensure it is valid YAML and try again.\n")
			os.Exit(1)
		}
	}

	nanofile.path = "./.nanofile"

	// then look for a local .nanofile and override any global, or remaining default
	// options found
	if fi, _ := os.Stat(nanofile.path); fi != nil {
		if err := ParseConfig(nanofile.path, nanofile); err != nil {
			fmt.Printf("Nanobox failed to parse your .nanofile. Please ensure it is valid YAML and try again.\n")
			os.Exit(1)
		}
	}

	// set name specific options after potential .nanofiles have been parsed
	nanofile.Domain = fmt.Sprintf("%s.nano.dev", nanofile.Name)
	nanofile.IP = appNameToIP(nanofile.Name)

	return nanofile
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
