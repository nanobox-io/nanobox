package share

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/jcelliott/lumber"

	"github.com/nanobox-io/nanobox/commands/server"
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util"
)

type Request struct {
	Path    string
	UID     int
	GID     int
	MountIP string
}

// EXPORTSFILE ...
var EXPORTSFILE = "/etc/exports"

func Exists(path string) bool {
	// read exports file
	existingFile, err := ioutil.ReadFile(EXPORTSFILE)
	if err != nil {
		// if i cant read the etc exports it doesnt exist
		return false
	}

	// get the provider because i need the mount ip
	provider, err := models.LoadProvider()
	if err != nil {
		return false
	}

	lineCheck := fmt.Sprintf("%s -alldirs -mapall=%v:%v", provider.MountIP, uid(), gid())

	lines := strings.Split(string(existingFile), "\n")

	for _, line := range lines {
		// get existing line
		if strings.Contains(line, lineCheck) {
			return strings.Contains(line, path+" ") || strings.Contains(line, path+"\" ")
		}
	}
	return false
}

func Add(path string) error {

	// get the provider because i need the mount ip
	provider, err := models.LoadProvider()
	if err != nil {
		return err
	}

	// create a rpc request
	req := Request{
		Path:    path,
		UID:     uid(),
		GID:     gid(),
		MountIP: provider.MountIP,
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
	err = server.ClientRun("ShareRPC.Add", req, resp)
	if err != nil || !resp.Success {
		err = fmt.Errorf("failed to add share %v %v", err, resp.Message)
	}
	return err
}

// the rpc function run from the server
func (sh *ShareRPC) Add(req Request, resp *Response) error {
	lumber.Info("req: %#v\n", req)

	// read exports file
	existingFile, err := ioutil.ReadFile(EXPORTSFILE)
	if err != nil {
		// if the file didnt exist lets create an empty existingFile
		existingFile = []byte("")
	}

	lineCheck := fmt.Sprintf("%s -alldirs -mapall=%v:%v", req.MountIP, req.UID, req.GID)

	lines := strings.Split(string(existingFile), "\n")

	found := false
	for i, line := range lines {
		// get existing line (mac exports are all on one line)
		if strings.Contains(line, lineCheck) {
			// add our path to the line
			// check to see if this path has already been added
			if !(strings.Contains(line, req.Path+" ") || strings.Contains(line, req.Path+"\" ")) {
				lines[i] = fmt.Sprintf("\"%s\" %s", req.Path, line)
			}

			lines[i] = cleanLine(lines[i], lineCheck)
			found = true
			break
		}
	}

	// if this is the first time nanobox adds to the exports (single line)
	if !found {
		lines = append(lines, fmt.Sprintf("\"%s\" %s", req.Path, lineCheck))
	}

	// save
	if err := ioutil.WriteFile(EXPORTSFILE, []byte(strings.Join(lines, "\n")), 0644); err != nil {
		return err
	}

	if err := reloadServer(); err != nil {
		return err
	}
	resp.Success = true
	return nil
}

func Remove(path string) error {

	// get the provider because i need the mount ip
	provider, err := models.LoadProvider()
	if err != nil {
		return err
	}

	// create a rpc request
	req := Request{
		Path:    path,
		UID:     uid(),
		GID:     gid(),
		MountIP: provider.MountIP,
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
	err = server.ClientRun("ShareRPC.Remove", req, resp)
	if err != nil || !resp.Success {
		err = fmt.Errorf("failed to remove share %v %v", err, resp.Message)
	}
	return err

}

// the rpc function run from the server
func (sh *ShareRPC) Remove(req Request, resp *Response) error {

	quotedPath := fmt.Sprintf("\"%s\"", req.Path)

	// read exports file
	existingFile, err := ioutil.ReadFile(EXPORTSFILE)
	if err != nil {
		// if the error exists the file didnt exist.
		lumber.Error("failed to read etc/exports: %s", err)
		return nil
	}

	lineCheck := fmt.Sprintf("%s -alldirs -mapall=%v:%v", req.MountIP, req.UID, req.GID)

	existingLines := strings.Split(string(existingFile), "\n")
	newLines := []string{}

	for _, line := range existingLines {
		// get existing line
		if !strings.Contains(line, lineCheck) {
			newLines = append(newLines, line)
			continue
		}

		// recreate the line without our path or quoted path
		line = strings.Replace(line, fmt.Sprintf("%s ", req.Path), "", 1)
		line = strings.Replace(line, fmt.Sprintf("%s ", quotedPath), "", 1)
		if line != lineCheck {
			// if there is still any paths left in our line
			line = cleanLine(line, lineCheck)
			newLines = append(newLines, line)
		}
	}

	// save
	if err := ioutil.WriteFile(EXPORTSFILE, []byte(strings.Join(newLines, "\n")), 0644); err != nil {
		return err
	}

	err = reloadServer()
	if err == nil {
		resp.Success = true
	}
	return err
}

// reloadServer will reload the nfs server with the new export configuration
func reloadServer() error {

	// dont reload the server when testing
	if flag.Lookup("test.v") != nil {
		return nil
	}

	if err := util.Retry(startNFSD, 5, time.Second); err != nil {
		lumber.Error("nfsd enable: %s", err)
		return err
	}

	// check the exports to make sure a reload will be successful; TODO: provide a
	// clear message for a direction to fix
	cmd := exec.Command("nfsd", "checkexports")
	if b, err := cmd.CombinedOutput(); err != nil {
		lumber.Debug("checkexports: %s", b)
		return fmt.Errorf("checkexports: %s %s", b, err.Error())
	}

	// update exports; TODO: provide a clear error message for a direction to fix
	cmd = exec.Command("nfsd", "update")
	if b, err := cmd.CombinedOutput(); err != nil {
		lumber.Debug("update: %s", b)
		return fmt.Errorf("update: %s %s", b, err.Error())
	}

	return nil
}

func startNFSD() error {
	// make sure nfsd is running
	cmd := exec.Command("nfsd", "enable")
	if b, err := cmd.CombinedOutput(); err != nil {
		lumber.Debug("enable nfs: %s", b)
		return fmt.Errorf("enable nfs: %s %s", b, err.Error())
	}

	// request it to start but if everythign is working starting could cause an error
	// we dont want to check
	exec.Command("nfsd", "start").CombinedOutput()

	// add a short delay because nfsd takes some time
	<-time.After(time.Second)

	// check to see if nfsd is running
	b, _ := exec.Command("netstat", "-ln").CombinedOutput()
	if !strings.Contains(string(b), ".111 ") {
		return fmt.Errorf("nfsd ports not in use")
	}
	return nil
}

func cleanLine(line, lineCheck string) string {
	// split on spaces and remove mount options. also cleanup stray quotes
	paths := strings.Split(strings.Replace(line, lineCheck, "", 1), "\" \"")
	paths[0] = strings.Replace(paths[0], "\"", "", 1)
	// the space in this `"\" "` is important. Prepending a space to lineCheck
	// after it gets passed in might work too, assuming it doesn't have one
	// already when getting passed in.
	paths[len(paths)-1] = strings.Replace(paths[len(paths)-1], "\" ", "", 1)

	goodPaths := []string{}
	for _, path := range paths {
		fileInfo, err := os.Stat(path)
		// continue on if the file doest exist or if it is not a directory
		if err != nil {
			lumber.Info("Failed to stat - %s", err.Error())
			continue
		}
		if !fileInfo.IsDir() {
			lumber.Info("Path is not a directory!")
			continue
		}

		goodPaths = append(goodPaths, path)
	}
	goodPaths = removeDuplicates(goodPaths)
	return fmt.Sprintf("\"%s\" %s", strings.Join(goodPaths, "\" \""), lineCheck)
}

// takes a set of paths and removes duplicates as well as cleaning up any child paths
func removeDuplicates(paths []string) []string {
	rtn := []string{}
	// look through the paths
	for i, path := range paths {
		// default to adding the path as a non duplicate
		add := true
		for j, originalPath := range paths {

			// if im looking at the same path then ignore it
			if i == j {
				continue
			}

			// todo: verify adding children on mac actually breaks (linux is ok)
			// todo: test parent, child, parent-er, child-er to ensure only the parent-most gets added
			// if i find an element that is shorter but the same directory structure
			// dont add the longer path
			if strings.HasPrefix(path, originalPath+"/") {
				lumber.Info("Parent directory already exported '%s'", path)
				add = false
			}
		}

		// if I didnt detect a shorter path then I need to add this one
		if add {
			rtn = append(rtn, path)
		}
	}
	return rtn
}
