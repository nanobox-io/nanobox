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
  configure     Configure Nanobox.
  run           Start your local development environment.
  build-runtime Build your app's runtime.
  compile-app   Compile your application.
  deploy        Deploy your application to a live remote or a dry-run environment.
  console       Open an interactive console inside a component.
  remote        Manage application remotes.
  status        Display the status of your Nanobox VM & apps.
  login         Authenticate your nanobox client with your nanobox.io account.
  logout        Remove your nanobox.io api token from your local nanobox client.
  clean         Clean out any apps that no longer exist.
  info          Show information about the specified environment.
  tunnel        Create a secure tunnel between your local machine & a live component.
  implode       Remove all Nanobox-created containers, files, & data.
  destroy       Destroy the current project and remove it from Nanobox.
  start         Start the Nanobox virtual machine.
  stop          Stop the Nanobox virtual machine.
  update-images Updates docker images.
  evar          Manage environment variables.
  dns           Manage dns aliases for local applications.
  log           Streams application logs.
  version       Show the current Nanobox version.
  server        Start a dedicated nanobox server

Flags:
      --debug     In the event of a failure, drop into debug context
  -h, --help      help for nanobox
  -t, --trace     Increases display output and sets level to trace
  -v, --verbose   Increases display output and sets level to debug

Use "nanobox [command] --help" for more information about a command.
```


### Documentation

- Nanobox documentation is available at [docs.nanobox.io](https://docs.nanobox.io/).
- Guides for popular languages, frameworks and services are avaialble at [guides.nanobox.io](http://guides.nanobox.io).


## Contributing

Contributing to Nanobox is easy. Just follow these [contribution guidelines](https://docs.nanobox.io/contributing/).
Nanobox uses [govendor](https://github.com/kardianos/govendor#the-vendor-tool-for-go) to vendor dependencies. Use `govendor sync` to restore dependencies.


### Contact

For help using Nanobox or if you have any questions/suggestions, please reach out to help@nanobox.io or find us on [slack](https://slack.nanoapp.io/). You can also [create a new issue on this project](https://github.com/nanobox-io/nanobox/issues/new).

[![nanobox logo](http://nano-assets.gopagoda.io/open-src/nanobox-open-src.png)](https://nanobox.io/open-source/)
