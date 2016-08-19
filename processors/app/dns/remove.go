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

// Remove ...
type Remove struct {
	App  models.App
	Name string
}

//
func (remove Remove) Run() error {

	// short-circuit if the entries dont exist; we do this after we validate meta
	// and load the app because both of those are needed to determin the entry we're
	// looking for
	if !remove.entriesExist() {
		return nil
	}

	// if we're not running as the privileged user, we need to re-exec with privilege
	if !util.IsPrivileged() {
		return remove.reExecPrivilege()
	}

	// remove the DNS entries
	return remove.removeEntries()
}

// entriesExist returns true if both entries already exist
func (remove *Remove) entriesExist() bool {

	// fetch the IPs
	envIP := remove.App.GlobalIPs["env"]

	// generate the entries
	env := dns.Entry(envIP, remove.Name, remove.App.ID)

	return dns.Exists(env)
}

// removeEntries removes the "dev" and "env" entries into the host dns
func (remove *Remove) removeEntries() error {

	// fetch the IPs
	envIP := remove.App.GlobalIPs["env"]

	// generate the entries
	env := dns.Entry(envIP, remove.Name, remove.App.ID)

	// remove the DNS entry from the /etc/hosts file
	if err := dns.Remove(env); err != nil {
		lumber.Error("dns:Remove:reExecPrivilege:dns.Remove(%s): %s", env, err)
		return err
	}

	return nil
}

// reExecPrivilege re-execs the current process with a privileged user
func (remove *Remove) reExecPrivilege() error {

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

	// call 'dev dns rm' with the original path and args; config.NanoboxPath() will be the
	// currently executing program, so this command will ultimately lead right back
	// here
	cmd := fmt.Sprintf("%s %s dns rm %s", config.NanoboxPath(), remove.App.Name, remove.Name)

	// if the sudo'ed subprocess fails, we need to return error to stop the process
	if err := util.PrivilegeExec(cmd); err != nil {
		lumber.Error("dns:Remove:reExecPrivilege:util.PrivilegeExec(%s): %s", cmd, err)
		return err
	}

	return nil
}
