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
3. Type `nanobox` and follow the instructions to authenticate the CLI with Nanobox.


### Usage
```go
Usage:
  nanobox [flags]
  nanobox [command]

Available Commands:
  run           Starts a nanobox, provisions the app, & runs the app's exec
  dev           Starts the nanobox, provisions app, & opens an interactive terminal
  info          Displays information about the nanobox and your app
  console       Opens an interactive terminal from inside your app on the nanobox
  destroy       Destroys the nanobox
  stop          Suspends the nanobox
  update        Updates the CLI to the newest available version
  update-images Updates the nanobox docker images
  box           Subcommands for managing the nanobox/boot2docker.box
  engine        Subcommands to aid in developing a custom engine

Flags:
      --background[=false]: Stops nanobox from auto-suspending.
  -f, --force[=false]: Forces a command to run (effects vary per command).
  -h, --help[=false]: help for nanobox
  -v, --verbose[=false]: Increase command output from 'info' to 'debug'.
      --version[=false]: Display the current version of this CLI

Additional help topics:
  nanobox production

Use "nanobox [command] --help" for more information about a command.
```

### Documentation

- Usage documentation is available at [nanobox.io](https://docs.nanobox.io/cli/).
- Source code documentation is available on [godoc](http://godoc.org/github.com/nanobox-io/nanobox-cli).


## Todo/Doing
- Tests!
- Make it work on Windows


## Contributing
1. Fork it
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Added some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create new Pull Request


### Contact

For help using the CLI or if you have any questions or suggestions, please find us on IRC (freenode) at #nanobox.
