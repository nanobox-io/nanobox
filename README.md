[![nanoboxdesktop logo](http://nano-assets.gopagoda.io/readme-headers/nanoboxdesktop.png)](http://nanobox.io/open-source#nanoboxdesktop)
[![Build Status](https://travis-ci.org/nanobox-io/nanobox.svg)](https://travis-ci.org/nanobox-io/nanobox)

## Nanobox

[Nanobox](https://nanobox.io/) automates the creation of isolated, repeatable environments for local and production applications. When developing locally, Nanobox provisions your app's infrastructure inside of a virtual machine (VM) and mounts your local codebase into the VM. Any changes made to your codebase are reflected inside the virtual environment.

Once code is built and tested locally, Nanobox provisions and deploys an identical infrastructure on a production platform.

## How It Works

Nanobox uses [Virtual Box](http://virtualbox.org) and [Docker](https://www.docker.com/) to create virtual development environments on your local machine. App configuration is handled in the [boxfile.yml](https://docs.nanobox.io/boxfile/), a small yaml config file used to provision and configure your apps' environments both locally and in production.

## Why Nanobox?

Nanobox allows you to stop configuring environments and just code. It guarantees that any project you start will work the same for anyone else collaborating on the project. When it's time to launch the project, you'll know that your production app will work, because it already works locally.

### Installation

By using the [Nanobox installer](https://nanobox.io/download). *(Recommended)* .The installer includes all required dependencies (Virtual Box & Docker).


### Usage
```
Usage:
  nanobox [flags]
  nanobox [command]

Available Commands:
  status      Displays the status of your Nanobox VM & running platforms.
  deploy      Deploys your generated build package to a production app.
  console     Opens an interactive console inside a production component.
  link        Manages links between local & production apps.
  login       Authenticates your nanobox client with your nanobox.io account.
  logout      Removes your nanobox.io api token from your local nanobox client.
  build       Generates a deployable build package.
  clean       Clean out any environemnts that no longer exist
  dev         Manages your 'development' environment.
  sim         Manages your 'simulated' environment.
  tunnel      Creates a secure tunnel between your local machine & a production component.
  destroy     Destroys the Nanobox virtual machine.
  start       Starts the Nanobox virtual machine.
  stop        Stop the Nanobox virtual machine.

Flags:
      --debug         Increases display output and sets level to debug
  -h, --help          help for nanobox
  -v, --verbose       Increases display output and sets level to debug
  -V, --veryverbose   Increases display output and sets level to trace

Use "nanobox [command] --help" for more information about a command.
```

### Documentation

- Nanobox documentation is available at [docs.nanobox.io](https://docs.nanobox.io/).
- Guides for popular languages, frameworks and services are avaialble at [guides.nanobox.io](http://guides.nanobox.io).


## Contributing
Contributing to Nanobox is easy. Just follow these [contribution guidelines](https://docs.nanobox.io/contributing/).

### Contact

For help using Nanobox or if you have any questions/suggestions, please reach out to help@nanobox.io or find us on IRC at #nanobox (freenode). You can also [create a new issue on this project](https://github.com/nanobox-io/nanobox/issues/new).

[![nanobox logo](http://nano-assets.gopagoda.io/open-src/nanobox-open-src.png)](http://nanobox.io/open-source)
