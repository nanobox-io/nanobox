package share

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/jcelliott/lumber"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util"
)

// EXPORTSFILE ...
var EXPORTSFILE = "/etc/exports"

func Exists(path string) bool {
	// open file
	b, err := ioutil.ReadFile(EXPORTSFILE)
	if err != nil {
		return false
	}
	// check to see if the path is in the file
	return bytes.Contains(b, []byte(path+" ")) || bytes.Contains(b, []byte(path+"\" "))
}

func Add(path string) error {

	// get the provider because i need the mount ip
	provider, err := models.LoadProvider()
	if err != nil {
		return err
	}

	// read exports file
	existingFile, err := ioutil.ReadFile(EXPORTSFILE)
	if err != nil {
		// if the file didnt exist lets create an empty existingFile
		existingFile = []byte("")
	}

	lineCheck := fmt.Sprintf("%s -alldirs -mapall=%v:%v", provider.MountIP, uid(), gid())

	lines := strings.Split(string(existingFile), "\n")

	found := false
	for i, line := range lines {
		// get existing line
		if strings.Contains(line, lineCheck) {
			// add our path to the line
			lines[i] = fmt.Sprintf("\"%s\" %s", path, line)
			lines[i] = cleanLine(lines[i], lineCheck)
			found = true
			break
		}
	}
	if !found {
		lines = append(lines, fmt.Sprintf("%s %s", path, lineCheck))
	}

	// save
	if err := ioutil.WriteFile(EXPORTSFILE, []byte(strings.Join(lines, "\n")), 0644); err != nil {
		return err
	}
	return reloadServer()
}

func Remove(path string) error {

	quotedPath := fmt.Sprintf("\"%s\"", path)
	// get the provider because i need the mount ip
	provider, err := models.LoadProvider()
	if err != nil {
		return err
	}

	// read exports file
	existingFile, err := ioutil.ReadFile(EXPORTSFILE)
	if err != nil {
		// if the error exists the file didnt exist.
		lumber.Error("failed to read etc/exports: %s", err)
		return nil
	}

	lineCheck := fmt.Sprintf("%s -alldirs -mapall=%v:%v", provider.MountIP, uid(), gid())

	existingLines := strings.Split(string(existingFile), "\n")
	newLines := []string{}

	for _, line := range existingLines {
		// get existing line
		if !strings.Contains(line, lineCheck) {
			newLines = append(newLines, line)
			continue
		}

		// recreate the line without our path or quoted path
		line = strings.Replace(line, fmt.Sprintf("%s ", path), "", 1)
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

	return reloadServer()
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

	// check to see if nfsd is running
	b, _ := exec.Command("netstat", "-ln").CombinedOutput()
	if !strings.Contains(string(b), ".111 ") {
		return fmt.Errorf("nfsd ports not in use")
	}
	return nil
}

func cleanLine(line, lineCheck string) string {
	paths := strings.Split(strings.Replace(line, lineCheck, "", 1), " ")
	goodPaths := []string{}
	for _, path := range paths {
		// remove the quotes from the path
		path = strings.Replace(path, "\"", "", -1)
		fileInfo, err := os.Stat(path)
		if err != nil || !fileInfo.IsDir() {
			// continue on if the file doest exist or if it is not a directory
			continue
		}
		goodPaths = append(goodPaths, path)
	}
	return fmt.Sprintf("\"%s\" %s", strings.Join(goodPaths, "\" \""), lineCheck)
}
