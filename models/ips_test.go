package models

import (
	"net"
  "testing"
)

func TestIPsSave(t *testing.T) {
  // clear the registry table when we're finished
  defer truncate("registry")
  
  ips := IPs{net.ParseIP("1.2.3.4")}
  
  err := ips.Save()
  if err != nil {
    t.Error(err)
  }
  
  // fetch the ips
  ips2 := IPs{}
  
  if err = get("registry", "ips", &ips2); err != nil {
    t.Errorf("failed to fetch ips: %s", err.Error())
  }
  
  if len(ips) != 1 {
    t.Errorf("ips doesn't match")
  }
}

func TestIPsDelete(t *testing.T) {
  // clear the registry table when we're finished
  defer truncate("registry")
  
	ips := IPs{net.ParseIP("1.2.3.4")}  

  if err := ips.Save(); err != nil {
    t.Error(err)
  }
  
  if err := ips.Delete(); err != nil {
    t.Error(err)
  }
  
  // make sure the auth is gone
  keys, err := keys("registry")
  if err != nil {
    t.Error(err)
  }
  
  if len(keys) > 0 {
    t.Errorf("auth was not deleted")
  }
}

func TestLoadIPs(t *testing.T) {
  // clear the registry table when we're finished
  defer truncate("registry")
  
  ips := IPs{net.ParseIP("1.2.3.4")}
  
  if err := ips.Save(); err != nil {
    t.Error(err)
  }
  
  ips2, err := LoadIPs()
  if err != nil {
    t.Error(err)
  }
  
  if len(ips2) != 1 {
    t.Errorf("did not load the correct ips")
  }
}
