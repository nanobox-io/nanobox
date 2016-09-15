// +build windows

package netfs

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/jcelliott/lumber"

	"github.com/nanobox-io/nanobox/util/display"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/display"
	"github.com/nanobox-io/nanobox/util/provider"
)

// Exists checks to see if the share already exists
func Exists(path string) bool {

	// running `net share` will list all of the cifs and other shares on the
	// windows machine. This can be run as a non-administrator.
	cmd := exec.Command("net", "share")
	output, err := cmd.CombinedOutput()

	// if there was an error, we'll short-circuit and return false
	if err != nil {
		return false
	}

	// return true if we find the path in the output
	if bytes.Contains(output, []byte(normalizePath(path))) {
		return true
	}

	return false
}

// Add exports a cifs share
func Add(path string) error {

	appID := config.EnvID()
	user := os.Getenv("USERNAME")

	// net share APPNAME=path /unlimited /GRANT:Everyone,FULL
	// net share APPNAME=path /unlimited /GRANT:%username%,FULL

	cmd := []string{
		"net",
		"share",
		fmt.Sprintf("nanobox-%s=%s", appID, normalizePath(path)),
		"/unlimited",
		fmt.Sprintf("/GRANT:%s,FULL", user),
	}

	process := exec.Command(cmd[0], cmd[1:]...)
	output, err := process.CombinedOutput()

	// if there was an error, we'll short-circuit and return false
	if err != nil {
		return err
	}

	// return nil (success) if the command was successful
	if bytes.Contains(output, []byte("was shared successfully.")) {
		return nil
	}

	// if we are here, it might have failed. Lets just check one more time
	// to see if the share already exists. If so, let's return success (nil)
	if Exists(path) {
		return nil
	}

	return errors.New("Failed to create cifs share")
}

// Remove removes a cifs share
func Remove(path string) error {

	appID := config.EnvID()

	// net share APPNAME /delete /y

	cmd := []string{
		"net",
		"share",
		fmt.Sprintf("nanobox-%s", appID),
		"/delete",
		"/y",
	}

	process := exec.Command(cmd[0], cmd[1:]...)
	output, err := process.CombinedOutput()

	// if there was an error, we'll short-circuit and return false
	if err != nil {
		return err
	}

	// return nil (success) if the command was successful
	if bytes.Contains(output, []byte("was deleted successfully.")) {
		return nil
	}

	// if we are here, it might have failed. Lets just check one more time
	// to see if the share is already gone. If so, let's return success (nil)
	if !Exists(path) {
		return nil
	}

	return errors.New("Failed to delete cifs share")
}

// Mount mounts a cifs share on a guest machine
func Mount(_, mountPath string) error {

	appID := config.EnvID()
	user := os.Getenv("USERNAME")

	// pause the current task
	display.PauseTask()
	// wait a bit to ensure the output doesn't get messed up
	<-time.After(time.Second * 1)
	
	// fetch the password from the user
	fmt.Printf("%s's password is required to mount a Windows share.\n", user)
	pass, err := display.ReadPassword()
	if err != nil {
		return err
	}
	
	// resume the task
	display.ResumeTask()

	// ensure the destination directory exists
	cmd := []string{"sudo", "/bin/mkdir", "-p", mountPath}
	if b, err := provider.Run(cmd); err != nil {
		lumber.Debug("output: %s", b)
		return fmt.Errorf("mkdir:%s", err.Error())
	}

	// ensure cifs/samba utilities are installed
	cmd = []string{"bash", "-c", setupCifsUtilsScript()}
	if b, err := provider.Run(cmd); err != nil {
		lumber.Debug("output: %s", b)
		return fmt.Errorf("mkdir:%s", err.Error())
	}

	// mount!
	// mount -t cifs -o username=USER,password=PASSWORD //192.168.99.1/APP /PATH
	source := fmt.Sprintf("//192.168.99.1/nanobox-%s", appID)
	opts := fmt.Sprintf("username=%s,password=%s,uid=1000,gid=1000", user, pass)
	cmd = []string{
		"sudo",
		"/bin/mount",
		"-t",
		"cifs",
		"-o",
		opts,
		source,
		mountPath,
	}
	if b, err := provider.Run(cmd); err != nil {
		lumber.Debug("output: %s", b)
		return fmt.Errorf("mount: output: %s err:%s", b, err.Error())
	}

	return nil
}

// normalizePath will ensure that a provided path is converted to a Windows
// style path including "\" instead of "/"
func normalizePath(path string) string {
	return strings.Replace(path, "/", "\\", -1)
}

// setupCifsUtilsScript returns a string containing the script to setup cifs
func setupCifsUtilsScript() string {
	script := `
		if [ ! -f /sbin/mount.cifs ]; then
			wget -O /mnt/sda1/tmp/tce/optional/samba-libs.tcz 
				http://repo.tinycorelinux.net/7.x/x86_64/tcz/samba-libs.tcz &&
			wget -O /mnt/sda1/tmp/tce/optional/cifs-utils.tcz 
				http://repo.tinycorelinux.net/7.x/x86_64/tcz/cifs-utils.tcz &&

			tce-load -i samba-libs &&
			tce-load -i cifs-utils
		fi
	`

	return strings.Replace(script, "\n", "", -1)
}
