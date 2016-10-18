package share

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/jcelliott/lumber"

	"github.com/nanobox-io/nanobox/util/config"
)

// Exists checks to see if the share already exists
func Exists(path string) bool {

	// running `net share` will list all of the cifs and other shares on the
	// windows machine. This can be run as a non-administrator.
	cmd := exec.Command("net", "share")
	output, err := cmd.CombinedOutput()
	lumber.Debug("net share output: %s", output)
	// if there was an error, we'll short-circuit and return false
	if err != nil {
		lumber.Debug("err from net share: %s", err)
		return false
	}

	// return true if we find the path in the output
	if bytes.Contains(output, []byte(path)) {
		lumber.Debug("output contains path: %s", path)
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
		fmt.Sprintf("nanobox-%s=%s", appID, path),
		"/unlimited",
		fmt.Sprintf("/GRANT:%s,FULL", user),
	}

	lumber.Debug("share add: %v", cmd)
	process := exec.Command(cmd[0], cmd[1:]...)
	output, err := process.CombinedOutput()
	lumber.Debug("share add output: %s", output)
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
	lumber.Debug("share remove: %v", cmd)
	process := exec.Command(cmd[0], cmd[1:]...)
	output, err := process.CombinedOutput()
	lumber.Debug("share remove output: %s", output)
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
