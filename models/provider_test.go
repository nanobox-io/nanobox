package models

import (
  "testing"
)

func TestProviderSave(t *testing.T) {
  // clear the registry table when we're finished
  defer truncate("registry")
  
  provider := Provider{
    HostIP: "192.168.1.2",
  }
  
  err := provider.Save()
  if err != nil {
    t.Error(err)
  }
  
  // fetch the provider
  provider2 := Provider{}
  
  if err = get("registry", "provider", &provider2); err != nil {
    t.Errorf("failed to fetch provider: %s", err.Error())
  }
  
  if provider2.HostIP != "192.168.1.2" {
    t.Errorf("provider doesn't match")
  }
}

func TestProviderDelete(t *testing.T) {
  // clear the registry table when we're finished
  defer truncate("registry")
  
  provider := Provider{
    HostIP: "192.168.1.2",
  }
  
  if err := provider.Save(); err != nil {
    t.Error(err)
  }
  
  if err := provider.Delete(); err != nil {
    t.Error(err)
  }
  
  // make sure the provider is gone
  keys, err := keys("registry")
  if err != nil {
    t.Error(err)
  }
  
  if len(keys) > 0 {
    t.Errorf("provider was not deleted")
  }
}

func TestLoadProvider(t *testing.T) {
  // clear the registry table when we're finished
  defer truncate("registry")
  
  provider := Provider{
    HostIP: "192.168.1.2",
  }
  
  if err := provider.Save(); err != nil {
    t.Error(err)
  }
  
  provider2, err := LoadProvider()
  if err != nil {
    t.Error(err)
  }
  
  if provider2.HostIP != "192.168.1.2" {
    t.Errorf("did not load the correct provider")
  }
}
