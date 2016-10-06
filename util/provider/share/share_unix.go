package share

import (
  "os"
  "strconv"
)

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
