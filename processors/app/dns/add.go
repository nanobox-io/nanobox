package dns

import (
	"fmt"
	"runtime"

	"github.com/jcelliott/lumber"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/dns"
)

// Add ...
type Add struct {
	App  models.App
	Name string
}

//
func (add Add) Run() error {

	// short-circuit if the entries already exist; we do this after we validate meta
	// and load the app because both of those are needed to determin the entry we're
	// looking for
	if add.entriesExist() {
		return nil
	}

	// if we're not running as the privileged user, we need to re-exec with privilege
	if !util.IsPrivileged() {
		return add.reExecPrivilege()
	}

	// add the DNS entries
	return add.addEntries()
}

// entriesExist returns true if both entries already exist
func (add *Add) entriesExist() bool {

	// fetch the IP
	// env in dev is used in the dev container
	// env in sim is used for portal
	envIP := add.App.GlobalIPs["env"]

	// generate the entries
	env := dns.Entry(envIP, add.Name, add.App.ID)

	// if the entries dont exist just return
	return dns.Exists(env)
}

// addEntries adds the dev and sim entries into the host dns
func (add *Add) addEntries() error {

	// fetch the IPs
	envIP := add.App.GlobalIPs["env"]

	// generate the entries
	env := dns.Entry(envIP, add.Name, add.App.ID)

	// add the 'sim' DNS entry into the /etc/hosts file
	err := dns.Add(env)
	if err != nil {
		lumber.Error("dns:Add:addEntries:dns.Add(%s): %s", env, err.Error())
		return fmt.Errorf("unalbe to add dns: %s", err.Error())
	}
	return nil
}

// reExecPrivilege re-execs the current process with a privileged user
func (add *Add) reExecPrivilege() error {

	if runtime.GOOS == "windows" {
		fmt.Println("Administrator privileges are required to modify host dns entries.")
		fmt.Println("Another window will be opened as the Administrator and permission may be requested.")

		// block here until the user hits enter. It's not ideal, but we need to make
		// sure they see the new window open if it requests their password.
		fmt.Println("Enter to continue:")
		var input string
		fmt.Scanln(&input)
	} else {
		fmt.Println("Root privileges are required to modify your hosts file, your password may be requested...")
	}

	// call 'dev dns add' with the original path and args; config.NanoboxPath() will be the
	// currently executing program with the path, so this command will ultimately lead right back
	// here
	cmd := fmt.Sprintf("%s %s dns add %s", config.NanoboxPath(), add.App.Name, add.Name)

	// if the sudo'ed subprocess fails, we need to return error to stop the process
	if err := util.PrivilegeExec(cmd); err != nil {
		lumber.Error("dns:Add:reExecPrivilege:util.PrivilegeExec(%s): %s", cmd, err.Error())
		return err
	}
	return nil
}
