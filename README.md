## Nanobox CLI

Releases available for the following Operating Systems (OS) and Architectures (OSArch):

Tested:

* OSX 10.6+ (darwin) - 386 | amd64
* Linux - 386 | amd64 | arm
* Windows 98+ - 386 | amd64

Untested:

* FreeBSD - 386 | amd64 | arm
* NetBSD - 386 | amd64 | arm
* OpenBSD - 386 | amd64
* Solaris - amd64

### Installation

To install the CLI first download the build that corresponds to your OS/OSArch from [nanobox.io/downloads](https://nanobox.io/downloads), then depending on your OS follow the steps below.


#### OSX and Linux:

1. Place the binary in your $PATH (ex. ~/bin) and run `chmod 755 nanobox`.
2. Type `nanobox` to view a list of available commands.


#### Windows

It is _**highly recommended**_ that you use [git bash](http://git-scm.com/downloads) as your shell. Not only will this provide the git workflow you'll need for Nanobox, but it will also give the CLI what it needs for features such as tunneling and SSH access to servers.

1. Place the nanobox.exe file in your $PATH (ex. C:/Windows), then you can either run it by clicking on it or using a shell.


#### From source:

1. `cd` into the folder on your GOPATH where you want the project to live.
2. `git clone git@github.com/nanobox-io/nanobox-cli`.
3. install [gom](https://github.com/mattn/gom). This is how the CLI manages dependancies. Skipping this step requires that you `go get` **each** dependancy in order to use the CLI.
3. `cd` into the new directory and run `gom install` or `gom build`.
4. Type `nanobox` and follow the instructions to authenticate the CLI with Nanobox.


### Usage

Typing `nanobox` will show the CLI's help text and list of available commands.

Typing `nanobox -v` or nanobox `--version` will tell you what version of the CLI you are currently using, and what the latest available version is.

All commands have a short (-) and a verbose (--) option when passing flags.

Passing `-h` to any command returns detailed information about that command.

Passing `--debug` at the end of any command will show all API request/response output.


#### Available Commands:

* `nanobox list` : List all your nanobox's.
* `nanobox info` : Display info about a nanobox.
* `nanobox create` : Provision a nanobox.
* `nanobox destroy` : Destroy a nanobox.
* `nanobox rebuild` : Rebuild and redeploy a nanobox.
* `nanobox rollback` : Roll a nanobox back one (1) deploy.
* `nanobox log` : Display nanobox log information.
* `nanobox open` : Open a nanobox in the default browser.
* `nanobox run` : Run a command on a nanobox service.
* `nanobox ssh` : Open an SSH connection to a nanobox service.
* `nanobox tunnel` : Create a port forward tunnel to a nanobox service.
* `nanobox evar:create` : Create an environment variable for nanobox
* `nanobox evar:destroy` : Destroy an environment variable for nanobox
* `nanobox evar:list` : List all environment variable for nanobox
* `nanobox service:list` : List an nanobox's services
* `nanobox service:info` : List info about a service
* `nanobox service:restart` : Restart a service
* `nanobox service:reboot` : Reboot a service
* `nanobox service:repair` : Repair a service


### Documentation

Complete documentation is available on [godoc](http://godoc.org/github.com/nanobox-io/nanobox-cli).


### Contact

For help using the CLI or if you have any questions or suggestions, please find us on IRC (freenode) at #nanobox.
