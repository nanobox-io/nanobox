package share

import (
	"fmt"
	"bytes"
	"io/ioutil"
	"strings"


	"github.com/nanobox-io/nanobox/models"
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
	return bytes.Contains(b, []byte(path + " "))
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
		return err
	}

	lineCheck := fmt.Sprintf("%s -alldirs -mapall=%v:%v", provider.MountIP, uid(), gid())

	lines := strings.Split(string(existingFile), "\n")

	found := false
	for i, line := range lines {
		// get existing line
		if strings.Contains(line, lineCheck) {
			// add our path to the line
			lines[i] = fmt.Sprintf("%s %s", path, line)
			found = true
			break
		}
	}
	if !found {
	  lines = append(lines, fmt.Sprintf("%s %s", path, lineCheck))
	}
	
	// save
	return ioutil.WriteFile(EXPORTSFILE, []byte(strings.Join(lines, "\n")), 0644)
}

func Remove(path string) error {
	// get the provider because i need the mount ip
	provider, err := models.LoadProvider()
	if err != nil {
		return err
	}

	// read exports file
	existingFile, err := ioutil.ReadFile(EXPORTSFILE)
	if err != nil {
		return err
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

		// add our path to the line
		line = strings.Replace(line, fmt.Sprintf("%s ", path), "", 1)
		if line != lineCheck {
			// if there is still any paths left in our line
			newLines = append(newLines, line)
		}
	}
	
	// save
	return ioutil.WriteFile(EXPORTSFILE, []byte(strings.Join(newLines, "\n")), 0644)
}