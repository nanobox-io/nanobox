package dns

import (
  "fmt"
  
  "github.com/jcelliott/lumber"
  
  "github.com/nanobox-io/nanobox/models"
  "github.com/nanobox-io/nanobox/util"
  "github.com/nanobox-io/nanobox/util/config"
  "github.com/nanobox-io/nanobox/util/display"
  "github.com/nanobox-io/nanobox/util/dns"
)

// Add adds a dns entry to the local hosts file
func Add(a *models.App, name string) error {
  
  // fetch the IP
  // env in dev is used in the dev container
  // env in sim is used for portal
  envIP := a.GlobalIPs["env"]
  
  // generate the dns entry
  entry := dns.Entry(envIP, name, a.ID)
  
  // short-circuit if this entry already exists
  if dns.Exists(entry) {
    return nil
  }
  
  // ensure we're running as the administrator for this
  if !util.IsPrivileged() {
    return reExecPrivilegedAdd(a, name)
  }
  
  // add the entry
  if err := dns.Add(entry); err != nil {
    lumber.Error("dns:Add:dns.Add(%s): %s", entry, err.Error())
    return fmt.Errorf("unable to add dns entry: %s", err.Error())
  }
  
  return nil
}

// reExecPrivilegedAdd re-execs the current process with a privileged user
func reExecPrivilegedAdd(a *models.App, name string) error {
  display.PrintRequiresPrivilege("to modify host dns entries")
  
  // call 'dev dns add' with the original path and args
  cmd := fmt.Sprintf("%s %s dns add %s", config.NanoboxPath(), a.Name, name)
  
  // if the sudo'ed subprocess fails, we need to return error to stop the process
  if err := util.PrivilegeExec(cmd); err != nil {
    lumber.Error("dns:reExecPrivilegedAdd:util.PrivilegeExec(%s): %s", cmd, err)
    return err
  }
  
  return nil
}
