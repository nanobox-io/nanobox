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

	"github.com/pagodabox/nanobox-boxfile"
	"github.com/pagodabox/nanobox-golang-stylish"
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
		Domain:   "gonano",
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
	fmt.Printf(stylish.Bullet("Parsing Boxfile"))

	//
	// if provider := nanobox.StringValue("provider"); provider != "" {
	//  fmt.Printf(stylish.Bullet(fmt.Sprintf("   - Custom Provider detected (%v)", provider)))
	// 	b.Provider = provider
	// }

	//
	if domain := nanobox.StringValue("domain"); domain != "" {
		fmt.Printf(stylish.Bullet(fmt.Sprintf("   - Custom domain detected (%v)", domain)))
		b.Domain = domain
	}

	//
	if ram := nanobox.IntValue("ram"); ram != 0 {

		// use specified RAM if it's greater than the default. Otherwise warn that
		// default will be used.
		switch {
		case ram >= b.RAM:
			fmt.Printf(stylish.Bullet(fmt.Sprintf("   - Custom RAM setting detected (%v)", ram)))
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
			fmt.Printf(stylish.Bullet(fmt.Sprintf("   - Custom CPU setting detected (%v)", cpus)))
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
			fmt.Printf(stylish.Bullet(fmt.Sprintf("   - Custom custom CPU cap detected (%v)", cpuCap)))
			b.CPUCap = cpuCap
		default:
			Console.Warn("Specified CPU cap (%v) is more than allowed default (%v), Using default...", cpuCap, b.CPUCap)
		}
	}

	//
	if engine := build.StringValue("engine"); engine != "" {
		fmt.Printf(stylish.Bullet(fmt.Sprintf("   - Engine detected (%v)", engine)))
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
