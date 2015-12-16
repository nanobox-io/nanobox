//
package vagrant

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/nanobox-io/nanobox-golang-stylish"

	"github.com/nanobox-io/nanobox/config"
	engineutil "github.com/nanobox-io/nanobox/util/engine"
)

// Init
func Init() {

	// create Vagrantfile
	vagrantfile, err := os.Create(config.AppDir + "/Vagrantfile")
	if err != nil {
		config.Fatal("[util/vagrant/init] ioutil.WriteFile() failed", err.Error())
	}
	defer vagrantfile.Close()

	//
	// create synced folders

	var sshMount, engineMount, codeMount string

	//
	// default path to ssh dir (assumes Unix)
	sshPath := filepath.Join(config.Home, ".ssh")

	// default path to ssh (windows)
	if config.OS == "windows" {
		sshPath = os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH") + `\.ssh`
	}

	// if an sshPath is provided in the nanofile override the default
	if config.Nanofile.SshPath != "" {
		sshPath = config.Nanofile.SshPath
	}

	// ensure the ssh location is a valid place; if the ssh dir is found, mount it
	sshDir, err := os.Stat(sshPath)
	if err == nil && sshDir.IsDir() {
		sshMount = fmt.Sprintf(`nanobox.vm.synced_folder '%s', "/mnt/ssh"`, sshPath)

		// if not found print this friendly warning
	} else {
		fmt.Printf(`
WARNING: Nanobox was unable to mount your .ssh folder into the VM because it was
unable to detect the location of an .ssh directory at:
%s

While nanobox is still usable for local development, this may result in failures
to fetch dependancies that require the use of those credentials`, sshPath)
	}

	//
	// mount code directory (mounted as nfs by default)
	codeMount = fmt.Sprintf(`nanobox.vm.synced_folder '%s', '/vagrant/code/%s'`, config.CWDir, config.Nanofile.Name)

	// mount code directory as NFS unless configured otherwise; if not mounted in
	// this way Vagrant will just decide what it thinks is best
	if config.Nanofile.MountNFS {
		codeMount += `,
      type: "nfs", mount_options: ["nfsvers=3", "proto=tcp"]`
	}

	//
	// "mount" the engine file locally at ~/.nanobox/apps/<app>/<engine>; this is
	// done when a local engine is detected so that the engine can be developed
	// and changed are reflected in the VM
	name, path, err := engineutil.MountLocal()
	if err != nil {
		config.Debug("No engine mounted (not found locally).")
	}

	// "mount" the engine into the VM (if there is one)
	if name != "" && path != "" {
		engineMount = fmt.Sprintf(`nanobox.vm.synced_folder '%s', "/vagrant/engines/%s"`, path, name)

		// mount engine directory as NFS unless configured otherwise; if not mounted
		// in this way Vagrant will just decide what it thinks is best
		if config.Nanofile.MountNFS {
			engineMount += `,
      type: "nfs", mount_options: ["nfsvers=3", "proto=tcp"]`
		}
	}

	//
	// nanofile config

	// create nanobox private network and unique forward port
	network := fmt.Sprintf("nanobox.vm.network \"private_network\", ip: \"%s\"", config.Nanofile.IP)
	sshport := fmt.Sprintf("nanobox.vm.network :forwarded_port, guest: 22, host: %v, id: 'ssh'", appNameToPort(config.Nanofile.Name))

	//
	provider := fmt.Sprintf(`# VirtualBox
    nanobox.vm.provider "virtualbox" do |p|
      p.name = "%v"

      p.customize ["modifyvm", :id, "--natdnshostresolver1", "%+v"]
      p.customize ["modifyvm", :id, "--cpuexecutioncap", "%v"]
      p.cpus = %v
      p.memory = %v
    end`, config.Nanofile.Name, config.Nanofile.HostDNS, config.Nanofile.CPUCap, config.Nanofile.CPUs, config.Nanofile.RAM)

	//
	// insert a provision script that will indicate to nanobox-server to boot into
	// 'devmode'
	var devmode string
	if config.Devmode {
		fmt.Printf(stylish.Bullet("Configuring vm to run in 'devmode'"))

		devmode = `# added because --dev was detected; boots the server into 'devmode'
    nanobox.vm.provision "shell", inline: <<-DEVMODE
      echo "Starting VM in dev mode..."
      mkdir -p /mnt/sda/var/nanobox
      touch /mnt/sda/var/nanobox/DEV
    DEVMODE`
	}

	//
	// insert a provision script that will allow utilization of system proxy vars
	var proxy string

	if config.Nanofile.UseProxy {
		fmt.Printf(stylish.Bullet("Configuring vm to use 'proxy' vars"))

		proxy = `# added because env 'use_proxy' was set; configures vm to use proxy
    nanobox.vm.provision "shell" do |s|
      s.inline = <<-PROXY
        echo "Configuring VM for proxy..."
        cat > /var/lib/boot2docker/profile <<-EOF
export http_proxy="$1"
export https_proxy="$2"
export https_user="$3"
export https_pass="$4"
EOF
        sudo /usr/local/etc/init.d/docker restart
      PROXY
      s.args = "'#{ENV['http_proxy']}' '#{ENV['https_proxy']}' '#{ENV['https_user']}' '#{ENV['https_pass']}'"
    end`
	}

	//
	// write to Vagrantfile
	vagrantfile.Write([]byte(fmt.Sprintf(`
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

#
Vagrant.configure(2) do |config|

  # add the boot2docker user credentials to allow nanobox to freely ssh into the vm
  # w/o requiring a password
  config.ssh.shell = "bash"
  config.ssh.username = "docker"
  config.ssh.password = "tcuser"

  config.vm.define :'%s' do |nanobox|

    ## Set the hostname of the vm to the app domain
    nanobox.vm.provision "shell", inline: <<-SCRIPT
      sudo hostname %s
    SCRIPT

    ## Wait for nanobox-server to be ready before vagrant exits
    nanobox.vm.provision "shell", inline: <<-WAIT
      echo "Waiting for nanobox server..."
      while ! nc -z 127.0.0.1 1757; do sleep 1; done;
    WAIT

    ## box
    nanobox.vm.box = "nanobox/boot2docker"


    ## network

    # add custom private network and ip and custom ssh port forward
    `+network+`
    `+sshport+`


    ## shared folders

    # disable default /vagrant share (overridden below)
    nanobox.vm.synced_folder ".", "/vagrant", disabled: true

    # add nanobox shared folders
    `+sshMount+`
    `+codeMount+`
    `+engineMount+`

    ## provider configs
    `+provider+`

    ## wait for the dhcp service to come online
    nanobox.vm.provision "shell", inline: <<-WAIT
      attempts=0
      while [[ ! -f /var/run/udhcpc.eth1.pid && $attempts -lt 30 ]]; do
        let attempts++
        sleep 1
      done
    WAIT

    # kill the eth1 dhcp server so that it doesn't override the assigned ip when
    # the lease is up
    nanobox.vm.provision "shell", inline: <<-KILL
      if [ -f /var/run/udhcpc.eth1.pid ]; then
        echo "Killing eth1 dhcp..."
        kill -9 $(cat /var/run/udhcpc.eth1.pid)
      fi
    KILL

    %s

    %s

  end
end`, config.Nanofile.Name, config.Nanofile.Domain, devmode, proxy)))
}

// appNameToPort generates a unique network port to allow running multiple vms at
// once
func appNameToPort(s string) string {

	port := 10000 // starting port is > than 100000 to try and avoid confilcts

	// create an md5 of the app name to ensure a uniqe port is generated each time
	h := md5.New()
	io.WriteString(h, s)

	// iterate through each byte in the md5 hash summing along the way
	for _, v := range []byte(h.Sum(nil)) {
		port += int(v)
	}

	return fmt.Sprint(port)
}
