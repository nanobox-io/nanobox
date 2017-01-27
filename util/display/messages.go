package display

import (
	"fmt"
	"os"
	"strings"
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
func InfoDevRunContainer(cmd, ip string) {
	os.Stderr.WriteString(fmt.Sprintf(`

      **
   *********
***************   Your command will run in an isolated Linux container
:: ********* ::   Code changes in either the container or desktop are mirrored
" ::: *** ::: "   ------------------------------------------------------------
  ""  :::  ""     If you run a server, access it at >> %s
    "" " ""
       "

RUNNING > %s
`, ip, cmd))

	os.Stderr.WriteString(fmt.Sprintf("%s\n", strings.Repeat("-", len(cmd)+10)))

}

func InfoSimDeploy(ip string) {
	os.Stderr.WriteString(fmt.Sprintf(`
--------------------------------------------------------------------------------
+ Your app is running in simulated production environment
+ Access your app at >> %s
--------------------------------------------------------------------------------

`, ip))
}

func DevRunEmpty() {
	os.Stderr.WriteString(fmt.Sprintf(`
! You don't have any web or worker start commands specified in your
  boxfile.yml. More information about start commands is available here:

  https://docs.nanobox.io/boxfile/web/#start-command

`))
}

func FirstDeploy() {
	os.Stderr.WriteString(fmt.Sprintf(`
--------------------------------------------------------------------------------
+ HEADS UP:
+ This is the first deploy to this app and the upload takes longer than usual.
+ Future deploys only sync the differences and will be much faster.
--------------------------------------------------------------------------------

`))
}

func FirstBuild() {
	os.Stderr.WriteString(fmt.Sprintf(`
--------------------------------------------------------------------------------
+ HEADS UP:
+ This is the first build for this project and will take longer than usual.
+ Future builds will pull from the cache and will be much faster.
--------------------------------------------------------------------------------

`))
}

func ProviderSetup() {
	os.Stderr.WriteString(fmt.Sprintf(`
--------------------------------------------------------------------------------
+ HEADS UP:
+ Nanobox will run a single VM transparently within VirtualBox.
+ All apps and containers will be launched within the same VM.
--------------------------------------------------------------------------------

`))
}

func MigrateOldRequired() {
	os.Stderr.WriteString(fmt.Sprintf(`
--------------------------------------------------------------------------------
+ WARNING:
+ Nanobox has been successfully upgraded! This change constitutes a major 
+ architectural refactor as well as data re-structure. To use this version we 
+ need to purge your current apps. No worries, nanobox will re-build them for 
+ you the next time you use "nanobox run".
--------------------------------------------------------------------------------
`))

}

func MigrateProviderRequired() {
	os.Stderr.WriteString(fmt.Sprintf(`
--------------------------------------------------------------------------------
+ WARNING:
+ It looks like you want to use a different provider, cool! Just FYI, we have
+ to bring down your existing apps as providers are not compatible. No worries, 
+ nanobox will re-build them for you the next time you use "nanobox run".
--------------------------------------------------------------------------------
`))
}

func BadTerminal() {
  os.Stderr.WriteString(fmt.Sprintf(`
--------------------------------------------------------------------------------
This console is currently not supported by nanobox
Please refer to the docs for more information
--------------------------------------------------------------------------------
`))
}

func MissingDependencies(provider string, missingParts []string) {
    fmt.Printf("Using nanobox with %s requires tools that appear to not be available on your system.\n", provider)
    fmt.Println(strings.Join(missingParts, "\n"))
    if provider == "native" {
      provider = "docker"
    }
    fmt.Printf("View these requirements at docs.nanobox.io/install/requirements/%s/\n", provider)
}