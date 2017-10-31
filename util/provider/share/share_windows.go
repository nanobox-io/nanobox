package share

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"

	"github.com/jcelliott/lumber"

	"github.com/nanobox-io/nanobox/commands/server"
	"github.com/nanobox-io/nanobox/util/config"
)

type Request struct {
	Path  string
	AppID string
	User  string
}

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

	// on windows we dont want to try adding more then once
	if Exists(path) {
		return nil
	}

	req := Request{
		Path:  path,
		AppID: config.EnvID(),
		User:  os.Getenv("USERNAME"),
	}
	resp := &Response{}

	// in testing we will call the rpc function directly
	if flag.Lookup("test.v") != nil {
		shareRPC := &ShareRPC{}
		err := shareRPC.Add(req, resp)
		if err != nil || !resp.Success {
			err = fmt.Errorf("failed to add share %v %v", err, resp.Message)
		}
		return err
	}

	// have the server run the share command
	err := server.ClientRun("ShareRPC.Add", req, resp)
	if err != nil || !resp.Success {
		err = fmt.Errorf("failed to add share %v %v", err, resp.Message)
	}
	return err

}

// the rpc function run from the server
func (sh *ShareRPC) Add(req Request, resp *Response) error {
	// net share APPNAME=path /unlimited /GRANT:%username%,FULL

	cmd := []string{
		"net",
		"share",
		fmt.Sprintf("nanobox-%s=%s", req.AppID, req.Path),
		"/unlimited",
		fmt.Sprintf("/GRANT:%s,FULL", req.User),
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
		resp.Success = true
		return nil
	}

	// if we are here, it might have failed. Lets just check one more time
	// to see if the share already exists. If so, let's return success (nil)
	if Exists(req.Path) {
		resp.Success = true
		return nil
	}

	return fmt.Errorf("Failed to create cifs share")
}

// Remove removes a cifs share
func Remove(path string) error {
	req := Request{
		Path:  path,
		AppID: config.EnvID(),
	}
	resp := &Response{}

	// in testing we will call the rpc function directly
	if flag.Lookup("test.v") != nil {
		shareRPC := &ShareRPC{}
		err := shareRPC.Remove(req, resp)
		if err != nil || !resp.Success {
			err = fmt.Errorf("failed to remove share %v %v", err, resp.Message)
		}
		return err
	}

	// have the server run the share command
	err := server.ClientRun("ShareRPC.Remove", req, resp)
	if err != nil || !resp.Success {
		err = fmt.Errorf("failed to remove share %v %v", err, resp.Message)
	}
	return err

}

// the rpc function run from the server
func (sh *ShareRPC) Remove(req Request, resp *Response) error {
	// net share APPNAME /delete /y

	cmd := []string{
		"net",
		"share",
		fmt.Sprintf("nanobox-%s", req.AppID),
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
		resp.Success = true
		return nil
	}

	// if we are here, it might have failed. Lets just check one more time
	// to see if the share is already gone. If so, let's return success (nil)
	if !Exists(req.Path) {
		resp.Success = true
		return nil
	}

	return fmt.Errorf("Failed to delete cifs share")

}
