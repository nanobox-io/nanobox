// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

package commands

//
import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/pagodabox/nanobox-cli/config"
	"github.com/pagodabox/nanobox-golang-stylish"
)

//
var initCmd = &cobra.Command{
	Hidden: true,

	Use:   "init",
	Short: "Creates a nanobox-flavored Vagrantfile",
	Long: `
Description:
  Creates a nanobox-flavored Vagrantfile

  -f, --force[=false]: Generate a fresh Vagrantfile (overriding the existing Vagrantfile)`,

	Run: nanoInit,
}

// nanoInit
func nanoInit(ccmd *cobra.Command, args []string) {

	//
	var provider, devmode string

	//
	// attempt to parse the boxfile first; we don't want to create an app folder
	// if the app isn't able to be created
	boxfile := config.ParseBoxfile()

	// creates a project folder at ~/.nanobox/apps/<name> (if it doesn't already
	// exists) where the Vagrantfile and .vagrant dir will live for each app
	if _, err := os.Stat(config.AppDir); err != nil {
		if err := os.Mkdir(config.AppDir, 0755); err != nil {
			config.Fatal("[commands/init] os.Mkdir() failed", err.Error())
		}
	}

	//
	// generate a Vagrantfile at ~/.nanobox/apps/<app-name>/Vagrantfile
	// only if one doesn't already exist (unless forced)
	if !fForce {
		if _, err := os.Stat(config.AppDir + "/Vagrantfile"); err == nil {
			return
		}
	}

	// create synced folders
	synced_folders := fmt.Sprintf("nanobox.vm.synced_folder \"%v\", \"/vagrant/code/%v\"", config.CWDir, config.Nanofile.Name)

	// if an engine path is provided, add it to the synced_folders
	if engine := boxfile.Build.Engine; engine != "" {
		if _, err := os.Stat(engine); err == nil {

			//
			fp, err := filepath.Abs(engine)
			if err != nil {
				config.Fatal("[commands/init] filepath.Abs() failed", err.Error())
			}

			base := filepath.Base(fp)

			synced_folders += fmt.Sprintf("\n    nanobox.vm.synced_folder \"%v\", \"/vagrant/engines/%v\"", fp, base)
		}
	}

	//
	// nanofile config
	//
	// create nanobox private network
	network := fmt.Sprintf("nanobox.vm.network \"private_network\", ip: \"%v\"", config.Nanofile.IP)

	//
	// configure provider
	switch config.Nanofile.Provider {

	//
	case "virtualbox":
		provider = fmt.Sprintf(`# VirtualBox
  nanobox.vm.provider "virtualbox" do |p|
    p.name = "%v"

    p.customize ["modifyvm", :id, "--cpuexecutioncap", "%v"]
    p.cpus = %v
    p.memory = %v
  end`, config.Nanofile.Name, config.Nanofile.CPUCap, config.Nanofile.CPUs, config.Nanofile.RAM)

	//
	case "vmware":
		provider = fmt.Sprintf(`# VMWare
  nanobox.vm.provider "vmware" do |p|
    v.vmx["numvcpus"] = "%v"
    v.vmx["memsize"] = "%v"
  end`, config.Nanofile.CPUCap, config.Nanofile.CPUs, config.Nanofile.RAM)
	}

	// command to pull the latest verison of boot2docker
	version := "`curl -s https://api.github.com/repos/pagodabox/nanobox-boot2docker/releases/latest | awk '/^  \"name\": / {print $2}' | tr -d ',\n\"'`.strip"

	// insert a provision script that will indicate to nanobox-server to boot into
	// 'devmode'
	if fDevmode {
		fmt.Printf(stylish.Bullet("Configuring vm to run in 'devmode'"))

		devmode = `# added because --dev was detected; boots the server into 'devmode'
    config.vm.provision "shell", inline: <<-DEVMODE
      echo "Starting VM in dev mode..."
      mkdir -p /mnt/sda/var/nanobox
      touch /mnt/sda/var/nanobox/DEV
    DEVMODE`
	}

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
## nanobox VM to fail! To regenerate this file, delete it and run             ##
## 'nanobox init'                                                             ##
##                                                                            ##
################################################################################

# -*- mode: ruby -*-
# vi: set ft=ruby :

# pull the latest version of nanobox-boot2docker
version = %v

#
Vagrant.configure(2) do |config|

	# add the boot2docker user credentials to allow nanobox to freely ssh into the vm
	# w/o requiring a password
	config.ssh.shell = "bash"
	config.ssh.username = "docker"
	config.ssh.password = "tcuser"

	config.vm.define :'%v' do |nanobox|

	  ## Wait for nanobox-server to be ready before vagrant exits
	  nanobox.vm.provision "shell", inline: <<-WAIT
      echo "Waiting for nanobox server..."
      while ! nc -z 127.0.0.1 1757; do sleep 1; done;
    WAIT

	  ## box
	  nanobox.vm.box_url = "https://github.com/pagodabox/nanobox-boot2docker/releases/download/#{version}/nanobox-boot2docker.box"
	  nanobox.vm.box     = "nanobox/boot2docker"


	  ## network
	  %s


	  ## shared folders

	  # disable default /vagrant share (overridden below)
	  nanobox.vm.synced_folder ".", "/vagrant", disabled: true

	  # add nanobox shared folders
	  %s


	  ## provider configs
	  %s

	  # kill the eth1 dhcp server so that it doesn't override the assigned ip when
	  # the lease is up
	  nanobox.vm.provision "shell", inline: <<-KILL
      echo "Killing eth1 dhcp..."
      kill -9 $(cat /var/run/udhcpc.eth1.pid)
    KILL

		%s

	end
end`, version, config.Nanofile.Name, network, synced_folders, provider, devmode)

	// write the Vagrantfile
	if err := ioutil.WriteFile(config.AppDir+"/Vagrantfile", []byte(vagrantfile), 0755); err != nil {
		config.Fatal("[commands/init] ioutil.WriteFile() failed", err.Error())
	}
}
