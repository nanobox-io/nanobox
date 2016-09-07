package evar

import (
  "fmt"
  
  "github.com/nanobox-io/nanobox/models"
  "github.com/nanobox-io/nanobox/util/display"
)

func Remove(appModel *models.App, keys []string) error {
  
  // delete the evars
  for _, key := range keys {
    delete(appModel.Evars, key)
  }
  
  // persist the app model
  if err := appModel.Save(); err != nil {
    return fmt.Errorf("failed to delete evars: %s", err.Error())
  }
  
  // print the deleted keys
  fmt.Println()
  for _, key := range keys {
    fmt.Printf("%s %s removed\n", display.TaskComplete, key)
  }
  fmt.Println()
  
  return nil
}
