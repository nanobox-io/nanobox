//
package config

import (
	"fmt"
	"github.com/nanobox-io/nanobox/util"
	"os"
	"path/filepath"
)

// NanofileConfig represents all available/expected .nanofile configurable options
type NanofileConfig struct {
	CPUCap   int    `json:"cpu_cap"`   // max %CPU usage allowed to the guest vm
	CPUs     int    `json:"cpus"`      // number of CPUs to dedicate to the guest vm
	Domain   string `json:"domain"`    // the domain to use in conjuntion with the ip when accesing the guest vm (defaults to <Name>.dev)
	IP       string `json:"ip"`        // the ip added to the /etc/hosts file for accessing the guest vm
	MountNFS bool   `json:"mount_nfs"` // does the code directory get mounted as NFS
	Name     string `json:"name"`      // the name given to the project (defaults to cwd)
	Provider string `json:"provider"`  // guest vm provider (virtual box, vmware, etc)
	RAM      int    `json:"ram"`       // ammount of RAM to dedicate to the guest vm
	HostDNS  string `json:"host_dns"`  // use the hosts dns resolver
}

// ParseNanofile
func ParseNanofile() NanofileConfig {

	//
	nanofile := NanofileConfig{
		CPUCap:   50,
		CPUs:     2,
		MountNFS: true,
		Name:     filepath.Base(CWDir),
		Provider: "virtualbox", // this may change in the future (adding additional hosts such as vmware)
		RAM:      1024,
		HostDNS:  "off",
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

	// set name specific options after potential .nanofiles have been parsed
	nanofile.Domain = fmt.Sprintf("%s.dev", nanofile.Name)

	// assign a default IP if none is specified
	if nanofile.IP == "" {
		nanofile.IP = util.StringToIP(nanofile.Name)
	}

	// if the OS is Windows folders CANNOT be mounted as NFS
	if OS == "windows" {
		nanofile.MountNFS = false
	}

	return nanofile
}
