package app

import (
  "net"

  "github.com/nanobox-io/nanobox/models"
  "github.com/nanobox-io/nanobox/processor"
  "github.com/nanobox-io/nanobox/util"
  "github.com/nanobox-io/nanobox/util/data"
  "github.com/nanobox-io/nanobox/util/ip_control"
)

type appTeardown struct {
  control processor.ProcessControl
  app     models.App
}

func init() {
  processor.Register("app_teardown", appTeardownFunc)
}

func appTeardownFunc(control processor.ProcessControl) (processor.Processor, error) {
  return &appTeardown{control: control}, nil
}

func (teardown *appTeardown) Results() processor.ProcessControl {
  return teardown.control
}

func (teardown *appTeardown) Process() error {

  if err := teardown.loadApp(); err != nil {
    return err
  }

  // short-circuit if the app isn't active
  if teardown.app.State == "initialized" {
    return nil
  }

  if err := teardown.releaseIPs(); err != nil {
    return err
  }

  if err := teardown.deleteApp(); err != nil {
    return err
  }

  return nil
}

// loadApp loads the app from the db
func (teardown *appTeardown) loadApp() error {
  // the app might not exist yet, so let's not return the error if it fails
  data.Get("apps", util.AppName(), &teardown.app)

  // set the default state
  if teardown.app.State == "" {
    teardown.app.State = "initialized"
  }

  return nil
}

// releaseIPs releases necessary app-global ip addresses
func (teardown *appTeardown) releaseIPs() error {

  if err := ip_control.ReturnIP(net.ParseIP(teardown.app.DevIP)); err != nil {
    return err
  }

  return nil
}

// deleteApp saves the app to the db
func (teardown *appTeardown) deleteApp() error {

  // delete the app model
  if err := data.Delete("apps", util.AppName()); err != nil {
    return err
  }

  return nil
}
