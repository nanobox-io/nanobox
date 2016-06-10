// +build darwin

package netfs

import (
  "os/exec"

  "github.com/jcelliott/lumber"
)

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
