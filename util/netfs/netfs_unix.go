// +build !windows

package netfs

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/jcelliott/lumber"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util/provider"
)

// EXPORTSFILE ...
const EXPORTSFILE = "/etc/exports"

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

// Mount mounts a share on a guest machine
func Mount(hostPath, mountPath string) error {

	// ensure portmap is running
	cmd := []string{"sudo", "/usr/local/sbin/portmap"}
	if b, err := provider.Run(cmd); err != nil {
		lumber.Debug("output: %s", b)
		return fmt.Errorf("portmap:%s", err.Error())
	}

	// ensure the destination directory exists
	cmd = []string{"sudo", "/bin/mkdir", "-p", mountPath}
	if b, err := provider.Run(cmd); err != nil {
		lumber.Debug("output: %s", b)
		return fmt.Errorf("mkdir:%s", err.Error())
	}

	// TODO: this IP shouldn't be hardcoded, needs to be figured out!
	source := fmt.Sprintf("192.168.99.1:%s", hostPath)
	cmd = []string{"sudo", "/bin/mount", "-t", "nfs", source, mountPath}
	if b, err := provider.Run(cmd); err != nil {
		lumber.Debug("output: %s", b)
		return fmt.Errorf("mount: output: %s err:%s", b, err.Error())
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

	entry := fmt.Sprintf("\"%s\" %s -alldirs -mapall=%v:%v", path, provider.MountIP, uid(), gid())

	return entry, nil
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
		if scanner.Text() == entry {
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

// uid will grab the original uid that called sudo if set
func uid() (uid int) {

	//
	uid = os.Geteuid()

	// if this process was started with sudo, sudo is nice enough to set
	// environment variables to inform us about the user that executed sudo
	//
	// let's see if this is the case
	if sudoUID := os.Getenv("SUDO_UID"); sudoUID != "" {

		// SUDO_UID was set, so we need to cast the string to an int
		if s, err := strconv.Atoi(sudoUID); err == nil {
			uid = s
		}
	}

	return
}

// gid will grab the original gid that called sudo if set
func gid() (gid int) {

	//
	gid = os.Getgid()

	// if this process was started with sudo, sudo is nice enough to set
	// environment variables to inform us about the user that executed sudo
	//
	// let's see if this is the case
	if sudoGid := os.Getenv("SUDO_GID"); sudoGid != "" {

		// SUDO_UID was set, so we need to cast the string to an int
		if s, err := strconv.Atoi(sudoGid); err == nil {
			gid = s
		}
	}

	return
}
