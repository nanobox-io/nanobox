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
	"strings"
)

// NanofileConfig represents all available/expected .nanofile configurable options
type NanofileConfig struct {
	CPUCap    int    `json:"cpu_cap"`    // max %CPU usage allowed to the guest vm
	CPUs      int    `json:"cpus"`       // number of CPUs to dedicate to the guest vm
	DevConfig string `json:"dev_config"` // the type of dev environment to configure on the guest vm
	Domain    string `json:"domain"`     // the domain to use in conjuntion with the ip when accesing the guest vm (defaults to <Name>.dev)
	HostDNS   string `json:"host_dns"`   // use the hosts dns resolver
	IP        string `json:"ip"`         // the ip added to the /etc/hosts file for accessing the guest vm
	MountNFS  bool   `json:"mount_nfs"`  // does the code directory get mounted as NFS
	Name      string `json:"name"`       // the name given to the project (defaults to cwd)
	Provider  string `json:"provider"`   // guest vm provider (virtual box, vmware, etc)
	RAM       int    `json:"ram"`        // ammount of RAM to dedicate to the guest vm
	SshPath   string `json:"ssh_path"`   // provide the path to the .ssh directory (if any)
	UseProxy  bool   `json:"use_proxy"`  // import http[s]_proxy variables into boot2docker
}

// ParseNanofile
func ParseNanofile() NanofileConfig {

	//
	nanofile := NanofileConfig{
		CPUCap:   50,
		CPUs:     2,
		HostDNS:  "off",
		MountNFS: true,
		Name:     filepath.Base(CWDir),
		Provider: "virtualbox", // this may change in the future (adding additional hosts such as vmware)
		RAM:      1024,
		UseProxy: false,
	}

	nanofilePath := Root + "/.nanofile"

	// look for a global .nanofile first in the ~/.nanobox directory, and override
	// any default options found.
	if _, err := os.Stat(nanofilePath); err == nil {
		if err := ParseConfig(nanofilePath, &nanofile); err != nil {
			fmt.Printf("Nanobox failed to parse your .nanofile. Please ensure it is valid YAML and try again.\n")
			Exit(1)
		}
	}

	nanofilePath = "./.nanofile"

	// then look for a local .nanofile and override any global, or remaining default
	// options found
	if _, err := os.Stat(nanofilePath); err == nil {
		if err := ParseConfig(nanofilePath, &nanofile); err != nil {
			fmt.Printf("Nanobox failed to parse your .nanofile. Please ensure it is valid YAML and try again.\n")
			Exit(1)
		}
	}

	// make sure the name doesnt have any spaces
	nanofile.Name = strings.Replace(nanofile.Name, " ", "-", -1)

	// set name specific options after potential .nanofiles have been parsed
	nanofile.Domain = fmt.Sprintf("%s.dev", nanofile.Name)

	// assign a default IP if none is specified
	if nanofile.IP == "" {
		nanofile.IP = appNameToIP(nanofile.Name)
	}

	// if the OS is Windows, folders CANNOT be mounted as NFS
	if OS == "windows" {
		nanofile.MountNFS = false
	}

	// if no dev config is provided, the default is "mount"
	if nanofile.DevConfig == "" {
		nanofile.DevConfig = "mount"
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

	// convert app name into a unique private network IP by adding the first portion
	// of the network with the generated portion
	binary.BigEndian.PutUint32(ip, (network + sum))

	return ip.String()
}
