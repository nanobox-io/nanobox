// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

package commands

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pagodabox/nanobox-cli/config"
	"github.com/pagodabox/nanobox-cli/ui"
	"github.com/pagodabox/nanobox-golang-stylish"
)

// InitCommand satisfies the Command interface
type InitCommand struct{}

// Help prints detailed help text for the app list command
func (c *InitCommand) Help() {
	ui.CPrint(`
Description:
  Creates a nanobox flavored Vagrantfile

Usage:
  nanobox init
  `)
}

// Run creates a Vagrantfile
func (c *InitCommand) Run(opts []string) {

	//
	// creates a project folder at ~/.nanobox/apps/<app-name> (if it doesn't already
	// exists) where the Vagrantfile and .vagrant dir will live for each app
	if di, _ := os.Stat(config.AppDir); di == nil {

		//
		fmt.Printf(stylish.Bullet("Creating project directory at: " + config.AppDir))

		if err := os.Mkdir(config.AppDir, 0755); err != nil {
			fmt.Println("There was an error creating a project directory for '%v' at '%v'. Exiting... %v", config.App, config.AppDir, err)
			os.Exit(1)
		}
	}

	//
	// generate a Vagrantfile at ~/.nanobox/apps/<app-name>/Vagrantfile if one doesn't
	// exist
	if fi, _ := os.Stat(config.AppDir + "/Vagrantfile"); fi == nil {

		//
		fmt.Printf(stylish.Bullet("Preparing nanobox Vagrantfile:"))
		fmt.Printf("   - Adding code directory mount at: /vagrant/code/%v\n", config.App)

		// create synced folders
		synced_folders := fmt.Sprintf("nanobox.vm.synced_folder \"%v\", \"/vagrant/code/%v\"", config.CWDir, config.App)

		// if an engine path is provided, add it to the synced_folders
		if engine := config.Boxfile.Engine; engine != "" {
			if fi, _ := os.Stat(engine); fi != nil {

				//
				fp, err := filepath.Abs(engine)
				if err != nil {
					ui.LogFatal("[commands.init] filepath.Abs() failed", err)
				}

				base := filepath.Base(fp)

				//
				fmt.Printf("   - Adding engine directory mount at: /vagrant/engines/%v\n", base)

				synced_folders += fmt.Sprintf("\n    nanobox.vm.synced_folder \"%v\", \"/vagrant/engines/%v\"", fp, base)
			}
		}

		//
		// create nanobox private network
		fmt.Printf("   - Adding nanobox private network (%v)\n", config.Boxfile.IP)
		network := fmt.Sprintf("nanobox.vm.network \"private_network\", ip: \"%v\"", config.Boxfile.IP)

		//
		// configure provider
		fmt.Printf("   - Adding detected provider (%v)\n", config.Boxfile.Provider)

		provider := ""

		//
		switch config.Boxfile.Provider {

		//
		case "virtualbox":
			provider = fmt.Sprintf(`# VirtualBox
    nanobox.vm.provider "virtualbox" do |p|
      p.name = "%v"

      p.customize ["modifyvm", :id, "--cpuexecutioncap", "%v"]
      p.cpus = %v
      p.memory = %v
    end`, config.App, config.Boxfile.CPUCap, config.Boxfile.CPUs, config.Boxfile.RAM)

		//
		case "vmware":
			provider = fmt.Sprintf(`# VMWare
    nanobox.vm.provider "vmware" do |p|
      v.vmx["numvcpus"] = "%v"
      v.vmx["memsize"] = "%v"
    end`, config.Boxfile.CPUCap, config.Boxfile.CPUs, config.Boxfile.RAM)
		}

		// command to pull the latest verison of boot2docker
		version := "`curl -s https://api.github.com/repos/pagodabox/nanobox-boot2docker/releases/latest | awk '/^  \"name\": / {print $2}' | tr -d ',\n\"'`.strip"

		//
		// create Vagrantfile
		vagrantfile := fmt.Sprintf(`
################################################################################
##                                                                            ##
##                                   ***                                      ##
##                                *********                                   ##
##                           *******************                              ##
##                       ***************************                          ##
##                           *******************                              ##
##                       ...      *********      ...                          ##
##                           ...     ***     ...                              ##
##                       +++      ...   ...      +++                          ##
##                           +++     ...     +++                              ##
##                       \\\      +++   +++      ///                          ##
##                           \\\     +++     ///                              ##
##                                \\     //                                   ##
##                                   \//                                      ##
##                                                                            ##
##                    _  _ ____ _  _ ____ ___  ____ _  _                      ##
##                    |\ | |__| |\ | |  | |__) |  |  \/                       ##
##                    | \| |  | | \| |__| |__) |__| _/\_                      ##
##                                                                            ##
## This file was generated by nanobox. Any modifications to it may cause your ##
## nanobox VM to fail! To regenerate this file run 'nanobox init'             ##
##                                                                            ##
################################################################################

# -*- mode: ruby -*-
# vi: set ft=ruby :

# pull the latest version of nanobox-boot2docker
version = %v

$wait = <<SCRIPT
echo "Waiting for nanobox server..."
while [ $(nc -z -w 4 127.0.0.1 1757) ]; do
  sleep 1
done
SCRIPT

Vagrant.configure(2) do |config|

  config.vm.define :nanobox_boot2docker do |nanobox|

    ## Wait for nanobox-server to be ready before vagrant exits
    nanobox.vm.provision "shell", inline: $wait


    ## box
    nanobox.vm.box_url = "https://github.com/pagodabox/nanobox-boot2docker/releases/download/#{version}/nanobox-boot2docker.box"
    nanobox.vm.box     = "nanobox/boot2docker"


    ## network
    %s


    ## shared folders

    # disable default /vagrant share to override...
    nanobox.vm.synced_folder ".", "/vagrant", disabled: true

    # ...add nanobox shared folders
    %s


    ## provider configs
    %s

  end

end`, version, network, synced_folders, provider)

		// write the Vagrantfile
		if err := ioutil.WriteFile(config.AppDir+"/Vagrantfile", []byte(vagrantfile), 0755); err != nil {
			ui.LogFatal("[commands.init] ioutil.WriteFile() failed", err)
		}

		//
		fmt.Println("   [√] nanobox Vagrantfile generated at: " + config.AppDir + "/Vagrantfile")
	}
}
