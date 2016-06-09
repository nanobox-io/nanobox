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

  "github.com/jcelliott/lumber"

  "github.com/nanobox-io/nanobox/util"
)

// Add will export an nfs share
func Add(host, path string) error {

  // This process requires root, check to see if we're the root user.
  // If not, we need to run a hidden command as sudo that will just call this
  // function again. Thus, the subprocess will be running as root
  // todo: sudo re-run thing
  if os.Geteuid() != 0 {
    // get the original nanobox executable
    nanobox := os.Args[0]

    // call dev netfs add with the original path (ultimately leads right back here)
    cmd := fmt.Sprintf("%s dev netfs add %s %s", nanobox, host, path)

    fmt.Println("Admin privileges are required to export an nfs share, your password may be requested...")

    // if the sudo'ed subprocess fails, we need to return error to stop the process
    if err := util.PrivilegeExec(cmd); err != nil {
      return err
    }

    // the subprocess exited successfully, so we can short-circuit here
    return nil
  }

  if !Exists(host, path) {
    // add entry into the /etc/exports file
    if err := addEntry(host, path); err != nil {
      return err
    }

    // reload nfsd
    if err := reloadServer(); err != nil {
      return err
    }
  }

  return nil
}

// Remove will remove an nfs share
func Remove(host, path string) error {

  // This process requires root, check to see if we're the root user.
  // If not, we need to run a hidden command as sudo that will just call this
  // function again. Thus, the subprocess will be running as root
  // todo: sudo re-run thing
  if os.Geteuid() != 0 {
    // get the original nanobox executable
    nanobox := os.Args[0]

    // call dev netfs add with the original path (ultimately leads right back here)
    cmd := fmt.Sprintf("%s dev netfs rm %s %s", nanobox, host, path)

    fmt.Println("Admin privileges are required to remove an nfs share, your password may be requested...")

    // if the sudo'ed subprocess fails, we need to return error to stop the process
    if err := util.PrivilegeExec(cmd); err != nil {
      return err
    }

    // the subprocess exited successfully, so we can short-circuit here
    return nil
  }

  if Exists(host, path) {
    // add entry into the /etc/exports file
    if err := removeEntry(host, path); err != nil {
      return err
    }

    // reload nfsd
    if err := reloadServer(); err != nil {
      return err
    }
  }

  return nil
}

// Exists checks to see if the mount already exists
func Exists(host, path string) bool {
	// open the /etc/exports file for scanning...
	f, err := os.Open("/etc/exports")
	if err != nil {
		return false
	}
	defer f.Close()

  // generate the exports entry
  entry, err := entry(host, path)
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

// Mount mounts a share on a guest machine
func Mount(host_path, mount_path string, context []string) error {
  // ensure portmap is running
  run := append(context, "/usr/local/sbin/portmap")
  cmd := exec.Command(run[0], run[1:]...)
  b, err := cmd.CombinedOutput()
  if err != nil {
    lumber.Debug("output: %s", b)
    return err
  }

  // ensure the destination directory exists
  run = append(context, []string{"/bin/mkdir", "-p", mount_path}...)
  cmd = exec.Command(run[0], run[1:]...)
  b, err = cmd.CombinedOutput()
  if err != nil {
    lumber.Debug("output: %s", b)
    return err
  }

  // mount!
  // todo: this IP shouldn't be hardcoded, needs to be figured out
  source := fmt.Sprintf("192.168.99.1:%s", host_path)
  run = append(context, []string{"/bin/mount", "-t", "nfs", source, mount_path}...)
  cmd = exec.Command(run[0], run[1:]...)
  b, err = cmd.CombinedOutput()
  if err != nil {
    lumber.Debug("output: %s", b)
    return err
  }

  return nil
}

// addEntry will add the entry into the /etc/exports file
func addEntry(host, path string) error {
  // open exports file
  f, err := os.OpenFile("/etc/exports", os.O_RDWR|os.O_APPEND, 0644)
  if err != nil {
    return err
  }
  defer f.Close()

  // generate the entry
  entry, err := entry(host, path)
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
func removeEntry(host, path string) error {
  var contents string

  // open exports file
  f, err := os.OpenFile("/etc/exports", os.O_RDWR, 0644)
  if err != nil {
    return err
  }
  defer f.Close()

  // generate the entry
  entry, err := entry(host, path)
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
func entry(host, path string) (string, error) {

  // fetch the uid/gid for the export statement
  uid := uid()
  gid := gid()

  message := "\"%s\" %s -alldirs -mapall=%v:%v"
  entry := fmt.Sprintf(message, path, host, uid, gid)

  return entry, nil
}

// reloadServer will reload the nfs server with the new export configuration
func reloadServer() error {
  // todo: make sure nfsd is enabled

  // check the exports to make sure a reload will be successful
  cmd := exec.Command("nfsd", "checkexports")
  b, err := cmd.CombinedOutput()
  if err != nil {
    // todo: provide a clear message for a direction to fix
    lumber.Debug("output: %s", b)
    return err
  }

  // update exports
  cmd = exec.Command("nfsd", "update")
  b, err = cmd.CombinedOutput()
  if err != nil {
    // todo: provide a clear message for a direction to fix
    lumber.Debug("output: %s", b)
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
