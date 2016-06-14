// +build !windows

package netfs

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/jcelliott/lumber"
)

// EXPORTSFILE ...
const EXPORTSFILE = "/etc/exports"

// Entry generates the mount entry for the exports file
func Entry(host, path string) string {
	return fmt.Sprintf("\"%s\" %s -alldirs -mapall=%v:%v", path, host, uid(), gid())
}

// Exists checks to see if the mount already exists
func Exists(entry string) bool {

	// open the /etc/exports file for scanning...
	f, err := os.Open(EXPORTSFILE)
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
func Add(entry string) error {

	// add entry into the /etc/exports file
	if err := addEntry(entry); err != nil {
		return err
	}

	// reload nfsd
	if err := reloadServer(); err != nil {
		return err
	}

	return nil
}

// Remove will remove an nfs share
func Remove(entry string) error {

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
func Mount(hostPath, mountPath string, context []string) error {

	// ensure portmap is running
	run := append(context, "/usr/local/sbin/portmap")
	cmd := exec.Command(run[0], run[1:]...)
	if b, err := cmd.CombinedOutput(); err != nil {
		lumber.Debug("output: %s", b)
		return err
	}

	// ensure the destination directory exists
	run = append(context, []string{"/bin/mkdir", "-p", mountPath}...)
	cmd = exec.Command(run[0], run[1:]...)
	if b, err := cmd.CombinedOutput(); err != nil {
		lumber.Debug("output: %s", b)
		return err
	}

	// TODO: this IP shouldn't be hardcoded, needs to be figured out mount!
	source := fmt.Sprintf("192.168.99.1:%s", hostPath)
	run = append(context, []string{"/bin/mount", "-t", "nfs", source, mountPath}...)
	cmd = exec.Command(run[0], run[1:]...)
	if b, err := cmd.CombinedOutput(); err != nil {
		lumber.Debug("output: %s", b)
		return err
	}

	return nil
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
