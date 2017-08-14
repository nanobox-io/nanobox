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

func TunnelEstablished(component, port string) {
	os.Stderr.WriteString(fmt.Sprintf(`
--------------------------------------------------------------------------------
+ Secure tunnel established to %s
+ Use the following credentials to connect 
--------------------------------------------------------------------------------

Host: 127.0.0.1
Port: %s
User: available in your dashboard (if applicable)
Pass: available in your dashboard (if applicable)

`, component, port))
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

func DeployComplete() {
	os.Stderr.WriteString(fmt.Sprintf(`
%s Success, this deploy is on the way!
  Check your dashboard for progress.

`, TaskComplete))
}

func LoginComplete() {
	os.Stderr.WriteString(fmt.Sprintf(`
%s You've successfully logged in
`, TaskComplete))
}

func NetworkCreateError(name, network string) {
	os.Stderr.WriteString(fmt.Sprintf(`
Nanobox is trying to create a native docker network, and it
looks like we have a conflict. An existing docker network is
already using the %s address space.

You will need to either remove the conflicting network, or set
an alternative address space with the following:

nanobox config set %s <unused ip/cidr>
`, network, name))
}

func VMCommunicationError() {
	os.Stderr.WriteString(fmt.Sprintf(`
--------------------------------------------------------------------------------
Nanobox has started a VM that needs access to your machine for mounting.
This VM is unable to communicate with the host machine currently. Please
verify that you don't have a firewall blocking this connection, and try again!
--------------------------------------------------------------------------------
`))
}

func NoGonanaUser() {
	os.Stderr.WriteString(fmt.Sprintf(`
We could not connect as the gonano user but we were able to 
fall back to the default user
`, TaskComplete))
}

func MissingBoxfile() {
	os.Stderr.WriteString(fmt.Sprintf(`
--------------------------------------------------------------------------------
MISSING BOXFILE.YML
Nanobox is looking for a boxfile.yml config file. You might want to 
check out our getting-started guide on configuring your app:

https://guides.nanobox.io/
--------------------------------------------------------------------------------
`))
}

func InvalidBoxfile() {
	os.Stderr.WriteString(fmt.Sprintf(`
--------------------------------------------------------------------------------
INVALID BOXFILE.YML
Nanobox requires valid yaml in your boxfile.yml config file. Please paste the
contents of your boxfile into www.yamllint.com to validate.
--------------------------------------------------------------------------------
`))
}

func TooManyKeys() {
	os.Stderr.WriteString(fmt.Sprintf(`
--------------------------------------------------------------------------------
POSSIBLY TOO MANY KEYS
Nanobox imports your ssh key directory for fetching dependencies but it appears
you may have more than we can handle. You might want to check out our docs on
specifying a key to use:

https://docs.nanobox.io/local-config/configure-nanobox/
--------------------------------------------------------------------------------
`))
}

func WorldWritable() {
	os.Stderr.WriteString(fmt.Sprintf(`
--------------------------------------------------------------------------------
Virtualbox was unable to create the virtual machine because a folder in the path
is globaly accessable and it should be private.
--------------------------------------------------------------------------------
`))

}

func LoginRequired() {
	os.Stderr.WriteString(fmt.Sprintf(`
It appears you are running Nanobox for the first time.
Login to your Nanobox account:
`))
}

func UnexpectedPrivilage() {
	os.Stderr.WriteString(fmt.Sprintf(`
--------------------------------------------------------------------------------
+ ERROR:
+ Nanobox is designed to run as a standard user (non root)
+ Please run all nanobox commands as a non privilage user
--------------------------------------------------------------------------------

`))
}

func BadPortType(protocol string) {
	os.Stderr.WriteString(fmt.Sprintf(`
--------------------------------------------------------------------------------
+ WARNING:
+ The boxfile.yml does not support port protocol '%s'. Using 'tcp' as default.
--------------------------------------------------------------------------------
`, protocol))
}

func PortInUse(port string) {
	os.Stderr.WriteString(fmt.Sprintf(`
--------------------------------------------------------------------------------
ADDRESS IN USE
It appears your local port (%s) is in use. Please specify a different port with
the '-p' flag. (eg. 'nanobox tunnel data.db -p 5444')
--------------------------------------------------------------------------------
`, port))
}

func ConsoleNodeNotFound() {
	os.Stderr.WriteString(fmt.Sprintf(`
--------------------------------------------------------------------------------
NODE NOT FOUND
It appears the node you are trying to console to does not exist. Please double
check your boxfile.yml. If your boxfile.yml does contain a node by the name you
specified, please contact support.
--------------------------------------------------------------------------------
`))
}
