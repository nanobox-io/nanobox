// +build darwin

package netfs

import (
  "os"
  "os/exec"
  "fmt"
  "bufio"
  "io/ioutil"
  "strconv"
  "strings"

  "github.com/nanobox-io/nanobox/provider"
  "github.com/nanobox-io/nanobox/util"
)

// Add will export an nfs share
func Add(path string) error {

  // This process requires root, check to see if we're the root user.
  // If not, we need to run a hidden command as sudo that will just call this
  // function again. Thus, the subprocess will be running as root
  // todo: sudo re-run thing
  if os.Geteuid() != 0 {
    // get the original nanobox executable
    nanobox := os.Args[0]

    // call dev netfs add with the original path (ultimately leads right back here)
    cmd := fmt.Sprintf("%s dev netfs add %s", nanobox, path)

    fmt.Println("Admin privileges are required to export an nfs share, your password may be requested...")

    // if the sudo'ed subprocess fails, we need to return error to stop the process
    if err := util.PrivilegeExec(cmd); err != nil {
      return err
    }

    // the subprocess exited successfully, so we can short-circuit here
    return nil
  }

  if !Exists(path) {
    // add entry into the /etc/exports file
    addEntry(path)

    // reload nfsd
    reloadServer()
  }

  return nil
}

// Remove will remove an nfs share
func Remove(path string) error {

  // This process requires root, check to see if we're the root user.
  // If not, we need to run a hidden command as sudo that will just call this
  // function again. Thus, the subprocess will be running as root
  // todo: sudo re-run thing
  if os.Geteuid() != 0 {
    // get the original nanobox executable
    nanobox := os.Args[0]

    // call dev netfs add with the original path (ultimately leads right back here)
    cmd := fmt.Sprintf("%s dev netfs rm %s", nanobox, path)

    fmt.Println("Admin privileges are required to remove an nfs share, your password may be requested...")

    // if the sudo'ed subprocess fails, we need to return error to stop the process
    if err := util.PrivilegeExec(cmd); err != nil {
      return err
    }

    // the subprocess exited successfully, so we can short-circuit here
    return nil
  }

  if Exists(path) {
    // add entry into the /etc/exports file
    removeEntry(path)

    // reload nfsd
    reloadServer()
  }

  return nil
}

// Exists checks to see if the mount already exists
func Exists(path string) bool {
	// open the /etc/exports file for scanning...
	f, err := os.Open("/etc/exports")
	if err != nil {
		return false
	}
	defer f.Close()

  // generate the exports entry
  entry, err := entry(path)
  if err != nil {
    return false
  }

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

// MountCmds returns a list of commands to mount the path on a linux guest
func MountCmds(path string) []string {
  // ensure portmap is running
  //  portmap
  // ensure the destination directory exists
  //  mkdir /mnt/dir
  // mount
  //  mount -t nfs ${server_ip}:${host_path} ${mount_path}
}

// addEntry will add the entry into the /etc/exports file
func addEntry(path string) error {
  // open exports file
  f, err := os.OpenFile("/etc/exports", os.O_RDWR|os.O_APPEND, 0644)
  if err != nil {
    return err
  }
  defer f.Close()

  // generate the entry
  entry, err := entry(path)
  if err != nil {
    return err
  }

  // write the entry to the file
  if _, err := f.WriteString(fmt.Sprintf("%s\n", entry)); err != nil {
    return err
  }

  return nil
}

// removeEntry will remove the entry from the /etc/exports file
func removeEntry(path string) error {
  var contents string

  // open exports file
  f, err := os.OpenFile("/etc/exports", os.O_RDWR, 0644)
  if err != nil {
    return err
  }
  defer f.Close()

  // generate the entry
  entry, err := entry(path)
  if err != nil {
    return err
  }

  // remove entry from /etc/hosts
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

  // write back the contents of the hosts file minus the removed entry
  if err := ioutil.WriteFile("/etc/exports", []byte(contents), 0644); err != nil {
    return err
  }

  return nil
}

// entry generates the mount entry for the exports file
func entry(path string) (string, error) {

  // we need to fetch the IP of the nanobox vm from the provider
  ip, err := provider.HostIP()
  if err != nil {
    return "", err
  }

  // fetch the uid/gid for the export statement
  uid := uid()
  gid := gid()

  message := "\"%s\" %s -alldirs -mapall=%v:%v"
  entry := fmt.Sprintf(message, path, ip, uid, gid)

  return entry, nil
}

// reloadServer will reload the nfs server with the new export configuration
func reloadServer() error {
  // todo: make sure nfsd is enabled


  // check the exports to make sure a reload will be successful
  cmd := exec.Command("nfsd checkexports")
  if err := cmd.Run(); err != nil {
    // todo: provide a clear message for a direction to fix
    return err
  }

  // update exports
  cmd = exec.Command("nfsd update")
  if err := cmd.Run(); err != nil {
    // todo: provide a clear message for a direction to fix
    return err
  }

  return nil
}

// uid will grab the original uid that called sudo if set
func uid() (uid int) {

  uid = os.Geteuid()

  // if this process was started with sudo, sudo is nice enough to set
  // environment variables to inform us about the user that executed sudo
  //
  // let's see if this is the case
  if sudoUid := os.Getenv("SUDO_UID"); sudoUid != "" {
    // SUDO_UID was set, so we need to cast the string to an int
    if s, err := strconv.Atoi(sudoUid); err == nil {
      uid = s
    }
  }

  return
}

// gid will grab the original gid that called sudo if set
func gid() (gid int) {

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
