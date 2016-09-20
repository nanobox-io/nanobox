package display

import (
	"fmt"
	"os"
)

func MOTD() {
	os.Stderr.WriteString(fmt.Sprintf(`
    
                                   **
                                ********
                             ***************
                          *********************
                            *****************
                          ::    *********    ::
                             ::    ***    ::
                           ++   :::   :::   ++
                              ++   :::   ++
                                 ++   ++
                                    +
                    _  _ ____ _  _ ____ ___  ____ _  _
                    |\ | |__| |\ | |  | |__) |  |  \/
                    | \| |  | | \| |__| |__) |__| _/\_
`))
}

func InfoProductionHost() {
	os.Stderr.WriteString(fmt.Sprintf(`
--------------------------------------------------------------------------------
+ WARNING:
+ You are on a live, production Linux server.
+ This host is primarily responsible for running docker containers.
+ Changes made to this machine have real consequences.
+ Proceed at your own risk.
--------------------------------------------------------------------------------

`))
}

func InfoProductionContainer() {
	os.Stderr.WriteString(fmt.Sprintf(`
--------------------------------------------------------------------------------
+ WARNING:
+ You are in a live, production Linux container.
+ Changes made to this machine have real consequences.
+ Proceed at your own risk.
--------------------------------------------------------------------------------

`))
}

func InfoLocalContainer() {
	os.Stderr.WriteString(fmt.Sprintf(`
--------------------------------------------------------------------------------
+ You are inside a Linux container on your local machine.
+ Anything here can be undone, so have fun and explore!
--------------------------------------------------------------------------------

`))
}

func InfoDevContainer(ip string) {
	os.Stderr.WriteString(fmt.Sprintf(`
--------------------------------------------------------------------------------
+ You are in a Linux container
+ Your local source code has been mounted into the container
+ Changes to your code in either the container or desktop will be mirrored
+ If you run a server, access it at >> %s
--------------------------------------------------------------------------------

`, ip))
}

func DevRunEmpty() {
	os.Stderr.WriteString(fmt.Sprintf(`
! You don't have any web or worker
  start commands specified in your
  boxfile.yml. More information about
  start commands is available here:
  
  docs.nanobox.io/app-config/boxfile/web/#start-command

`))
}

func FirstDeploy() {
	os.Stderr.WriteString(fmt.Sprintf(`
--------------------------------------------------------------------------------
+ HEADS UP:
+ This is the first deploy to this app and the upload takes longer than usual.
+ Future deploys only upload differences and will be much faster.
--------------------------------------------------------------------------------

`))
}
