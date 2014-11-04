package main

import (
	"github.com/nanobox-core/cli/ui"
)

// Help prints out the help text for the entire CLI
func (cli *CLI) Help() {
	ui.CPrintln(`
                                    [white]/^[reset]
[cyan]. . . . . . . . . . . . . . . . .[reset] [white]////\\[reset] [cyan]. . . . . . . . . . . . . . . . . . . .[reset]
[blue]                                /  / /   \
                               / \         \
                              /    \       /\
                             / \ /   \   /    \
                             \  |      /      /
                               \|      \  / /
                                 \      |//
                                   \  //
                                     ' [reset]
    _____   _____   ______  _____   ____   _____       ______  __      _____
   |  .  | |  .  | |   __/ |     | |    \ |  .  |     |   __/ |  |    |_   _|
   |   __| |     | |  |  | |  .  | |  .  ||     |  _  |  |__  |  |__   _| |_
   |__|    |__|__| |_____| |_____| |____/ |__|__| |_| |_____| |_____| |_____|


Description:
  Welcome to the Nanobox CLI! We hope this tool will help your workflow when
  using Nanobox from the command line. If you encounter any issues or have
  any suggestions, [green]find us on IRC (freenode) at #pagodabox[reset]. Our
  engineers are available between 8 - 5pm MST.

  All commands have a short [-*] and a verbose [--*] option when passing flags.

  You can pass -h, --help, or help to any command to receive detailed information
  about that command.

  You can pass --debug at the end of any command to see all request/response
  output when making API calls.

Usage:
  pagoda (<COMMAND>:<ACTION> OR <ALIAS>) [GLOBAL FLAG] <POSITIONAL> [SUB FLAGS] [--debug]

Options:
  -h, --help, help
    Run anytime to receive detailed information about a command.

  -v, --version, version
    Run anytime to see the current version of the CLI.

  --debug
    Shows all API request/response output. [red]MUST APPEAR LAST[reset]

Available Commands:

  user            : Display your users information.

  list            : List all your applications.
  info            : Display info about an application.
  create          : Create an application.
  destroy         : Destroy an app.
  rebuild         : Rebuild and redeploy an app.
  rollback        : Roll an app back one (1) deploy.
  log             : Display app log information.
  open            : Open an app in the default browser.

  run             : Run a command on an app service.
  ssh             : Open an SSH connection to an app service.
  tunnel          : Create a port forward tunnel to an app service.

  evar:create     : Create an environment variable for app
  evar:destroy    : Destroy an environment variable for app
  evar:list       : List all environment variable for app

  service:list    : List an app's services
  service:info    : List info about a service
  service:restart : Restart a service
  service:reboot  : Reboot a service
  service:repair  : Repair a service
  `)
}
