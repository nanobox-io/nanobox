package share

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/jcelliott/lumber"

	"github.com/nanobox-io/nanobox/models"
)

// EXPORTSFILE ...
var EXPORTSFILE = "/etc/exports"

// Exists checks to see if the mount already exists
func Exists(path string) bool {

	// generate the entry
	entry, err := entry(path)
	if err != nil {
		return false
	}

	// open the /etc/exports file for scanning...
	var f *os.File
	f, err = os.Open(EXPORTSFILE)
	if err != nil {
		return false
	}
	defer f.Close()

	// scan exports file looking for an entry for this path...
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		// scan each line to see if we have a match​
		if scanner.Text() == entry {
			return true
		}
	}

	return false
}

// Add will export an nfs share
func Add(path string) error {

	// generate the entry
	entry, err := entry(path)
	if err != nil {
		return err
	}

	// add entry into the /etc/exports file
	if err := addEntry(entry); err != nil {
		return err
	}

	if err := cleanExport(); err != nil {
		return err
	}

	// reload nfsd
	if err := reloadServer(); err != nil {
		return err
	}

	return nil
}

// Remove will remove an nfs share
func Remove(path string) error {

	// generate the entry
	entry, err := entry(path)
	if err != nil {
		return err
	}

	if err := removeEntry(entry); err != nil {
		return err
	}

	// reload nfsd
	if err := reloadServer(); err != nil {
		return err
	}

	return nil
}

// entry generates the mount entry for the exports file
func entry(path string) (string, error) {

	// use the mountIP saved on the provider in the database
	provider, err := models.LoadProvider()
	if err != nil {
		return "", err
	}

	if provider.MountIP == "" {
		return "", fmt.Errorf("there is no mount ip on the provider")
	}

	return fmt.Sprintf("%s %s(rw,crossmnt,sync,no_subtree_check,all_squash,anonuid=%d,anongid=%d)", path, provider.MountIP, uid(), gid()), nil
}

// addEntry will add the entry into the /etc/exports file
func addEntry(entry string) error {

	// open exports file
	f, err := os.OpenFile(EXPORTSFILE, os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	// write the entry to the file
	if _, err := f.WriteString(fmt.Sprintf("%s\n", entry)); err != nil {
		return err
	}

	return nil
}

func cleanExport() error {

	// contents will end up storing the entire contents of the file excluding the
	// entry that no longer have a folder
	var contents string

	// open exports file
	f, err := os.OpenFile(EXPORTSFILE, os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	// remove entry from /etc/exports
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), "\"")

		// if it starts with a "
		if len(parts) > 1 {
			fileInfo, err := os.Stat(parts[1])
			if err != nil || !fileInfo.IsDir() {
				continue
			}
		}

		// add each line back into the file
		contents += fmt.Sprintf("%s\n", scanner.Text())
	}

	// trim the contents to avoid any extra newlines
	contents = strings.TrimSpace(contents)

	// add a single newline for completeness
	contents += "\n"

	// write back the contents of the exports file minus the removed entry
	if err := ioutil.WriteFile(EXPORTSFILE, []byte(contents), 0644); err != nil {
		return err
	}

	return nil
}

// removeEntry will remove the entry from the /etc/exports file
func removeEntry(entry string) error {

	// contents will end up storing the entire contents of the file excluding the
	// entry that is trying to be removed
	var contents string

	// open exports file
	f, err := os.OpenFile(EXPORTSFILE, os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	// remove entry from /etc/exports
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {

		// if the line contain the entry skip it
		if strings.Contains(scanner.Text(), entry) {
			continue
		}

		// add each line back into the file
		contents += fmt.Sprintf("%s\n", scanner.Text())
	}

	// trim the contents to avoid any extra newlines
	contents = strings.TrimSpace(contents)

	// add a single newline for completeness
	contents += "\n"

	// write back the contents of the exports file minus the removed entry
	if err := ioutil.WriteFile(EXPORTSFILE, []byte(contents), 0644); err != nil {
		return err
	}

	return nil
}

// reloadServer reloads the nfs server with the new export configuration
func reloadServer() error {
	// dont reload the server when testing
	if flag.Lookup("test.v") != nil {
		return nil
	}

	// make sure nfsd is running
	cmd := exec.Command("service", "nfs-server", "start")
	if b, err := cmd.CombinedOutput(); err != nil {
		lumber.Debug("enable nfs: %s err: %s", b, err)
	}

	// reload nfs server
	//  TODO: provide a clear error message for a direction to fix
	cmd := exec.Command("exportfs", "-ra")
	if b, err := cmd.CombinedOutput(); err != nil {
		lumber.Debug("update: %s", b)
		return fmt.Errorf("update: %s %s", b, err.Error())
	}

	return nil
}
