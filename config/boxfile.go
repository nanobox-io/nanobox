// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

package config

import (
	"encoding/binary"
	"net"

	"github.com/pagodabox/nanobox-boxfile"
)

// BoxfileConfig represents all available/expected Boxfile configurable options
type BoxfileConfig struct {

	// nanobox specific
	CPUCap   int
	CPUs     int
	Domain   string
	IP       string
	RAM      int
	Provider string

	// build
	Engine string
}

// ParseBoxfile
func ParseBoxfile() *BoxfileConfig {

	// default boxfile config options
	b := &BoxfileConfig{

		// nanobox
		IP:       appNameToIP(App),
		Domain:   "gonano.dev",
		Provider: "virtualbox",
		CPUCap:   50,
		CPUs:     2,
		RAM:      512,

		// build options
		Engine: "",
	}

	//
	f := boxfile.NewFromPath(CWDir + "/" + "Boxfile")
	nanobox := f.Node("nanobox")
	build := f.Node("build")

	//
	Console.Info("Parsing Boxfile...")

	//
	// if provider := nanobox.StringValue("provider"); provider != "" {
	// 	b.Provider = provider
	// }

	//
	if domain := nanobox.StringValue("domain"); domain != "" {
		Console.Info("Custom domain detected (%v), overriding default (%v)...", domain, b.Domain)
		b.Domain = domain
	}

	//
	if ram := nanobox.IntValue("ram"); ram != 0 {

		// use specified RAM if it's greater than the default. Otherwise warn that
		// default will be used.
		switch {
		case ram >= b.RAM:
			Console.Info("Using custom RAM setting (%v)...", ram)
			b.RAM = ram
		default:
			Console.Warn("Specified RAM (%v) is less than allowed default (%v), Using default...", ram, b.RAM)
		}
	}

	//
	if cpus := nanobox.IntValue("cpus"); cpus != 0 {

		// use specified CPUs if it's greater than the default. Otherwise warn that
		// default will be used.
		switch {
		case cpus >= b.CPUs:
			Console.Info("Using custom CPUs setting (%v)...", cpus)
			b.CPUs = cpus
		default:
			Console.Warn("Specified CPUs (%v) is less than allowed default (%v), Using default...", cpus, b.CPUs)
		}
	}

	//
	if cpuCap := nanobox.IntValue("cpu_cap"); cpuCap != 0 {

		// use specified CPU cap if it's greater than the default. Otherwise warn that
		// default will be used.
		switch {
		case cpuCap <= b.CPUCap:
			Console.Info("Using custom CPU Cap (%v)...", cpuCap)
			b.CPUCap = cpuCap
		default:
			Console.Warn("Specified CPU cap (%v) is more than allowed default (%v), Using default...", cpuCap, b.CPUCap)
		}
	}

	//
	if engine := build.StringValue("engine"); engine != "" {
		Console.Info("Engine detected (%v)...", engine)
		b.Engine = engine
	}

	return b

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
