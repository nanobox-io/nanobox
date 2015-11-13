[![nanoboxdesktop logo](http://nano-assets.gopagoda.io/readme-headers/nanoboxdesktop.png)](http://nanobox.io/open-source#nanoboxdesktop)
[![Build Status](https://travis-ci.org/nanopack/nanobox.svg)](https://travis-ci.org/nanopack/nanobox)

## Nanobox Desktop

[Nanobox desktop](https://desktop.nanobox.io/) gets rid of conventional "environments", such as development and production, and creates one unified environment; a nanobox environment. [Engines](https://docs.nanobox.io/engines/) help nanobox create these isolated, shareable, disposable environments, ensuring that whatever environment is created with nanobox is the same no matter where it lives.


## How It Works

Leveraging Vagrant and Docker, nanobox creates a virtual machine and launches docker containers that represent each piece of your application. Your nanobox environment will have everything it needs to run your application. Once you are done with it, you can "throw it away", leaving your machine clean. You'll never have to worry about having the right version of a language installed or managing dependencies again.


## Why Nanobox?

Nanobox allows you to stop configuring environments and just code. It guarantees that any project you start will work the same for anyone else collaborating on the project. When it's time to launch the project, you'll know that your production app will work, because it already works on nanobox.


### Installation

[Nanobox desktop](https://desktop.nanobox.io/downloads/) can be installed in two different ways:

1. Just the nanobox desktop binary.
2. The nanobox desktop installer which includes the binary and all its dependencies (Vagrant, and Virtualbox)


#### When using just the binary:

##### OSX and Linux:

1. Place the binary in your $PATH (ex. ~/bin) and run "chmod 755 nanobox".
3. Run "which nanobox" to ensure it's accessible from your $PATH.
2. Type `nanobox` to view a list of available commands.


##### Windows

It is _**highly recommended**_ that you use [git bash](http://git-scm.com/downloads) as your shell. Not only will this provide the git workflow you'll need for Nanobox, but it will also give the CLI what it needs for features such as tunneling and SSH access to servers.

1. Place the nanobox.exe file in your $PATH (ex. C:/Windows), then you can either run it by clicking on it or using a shell.


### Usage
```
Usage:
  nanobox [flags]
  nanobox [command]

Available Commands:
  run           Starts a nanobox, provisions the app, & runs the app's exec
  dev           Starts the nanobox, provisions app, & opens an interactive terminal
  info          Displays information about the nanobox and your app
  console       Opens an interactive terminal from inside your app on nanobox
  destroy       Destroys the nanobox
  stop          Suspends the nanobox
  update        Updates the CLI to the newest available version
  update-images Updates the nanobox docker images
  box           Subcommands for managing the nanobox/boot2docker.box
  engine        Subcommands to aid in developing a custom engine

Flags:
      --background[=false]: Stops nanobox from auto-suspending.
  -f, --force[=false]: Forces a command to run (effects vary per command).
  -v, --verbose[=false]: Increase command output from 'info' to 'debug'.
      --version[=false]: Display the current version of this CLI

Additional help topics:
  nanobox production

Use "nanobox [command] --help" for more information about a command.
```

### Documentation

- Usage documentation is available at [nanobox.io](https://docs.nanobox.io/cli/).
- Source code documentation is available on [godoc](http://godoc.org/github.com/nanobox-io/nanobox).


## Todo/Doing
- Make it work on Windows
- More tests (you can never had enough!)


## Contributing
Contributing to nanobox desktop is easy, just follow this [guide](https://docs.nanobox.io/contributing/)


### Contact

For help using nanobox desktop or if you have any questions/suggestions, please find us on IRC (freenode) at #nanobox.

[![nanobox logo](http://nano-assets.gopagoda.io/open-src/nanobox-open-src.png)](http://nanobox.io/open-source)
