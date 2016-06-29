package env

import (
  "fmt"

  "github.com/nanobox-io/nanobox/processor"
  "github.com/nanobox-io/nanobox/provider"
  "github.com/nanobox-io/nanobox/util/config"
  "github.com/nanobox-io/nanobox/util/netfs"
)

// processEnvSetup ...
type processEnvSetup struct {
  control processor.ProcessControl
}

//
func init() {
  processor.Register("env_setup", envSetupFn)
}

//
func envSetupFn(control processor.ProcessControl) (processor.Processor, error) {
  // control.Meta["processEnvSetup-control"]

  // do some control validation check on the meta for the flags and make sure they
  // work

  return &processEnvSetup{control: control}, nil
}

//
func (envSetup processEnvSetup) Results() processor.ProcessControl {
  return envSetup.control
}

//
func (envSetup *processEnvSetup) Process() error {

  if err := envSetup.setupProvider(); err != nil {
    return err
  }

  if err := envSetup.setupMounts(); err != nil {
    return err
  }

  // if there is an environment then we should set up app
  // if not (in the case of a build) no app setup is necessary
  if envSetup.control.Env != "" {
    if err := envSetup.setupApp(); err != nil {
      return err
    }
  }

  return nil
}

// setupProvider sets up the provider
func (envSetup *processEnvSetup) setupProvider() error {
  return processor.Run("provider_setup", envSetup.control)
}

// setupMounts will add the envs and mounts for this app
func (envSetup *processEnvSetup) setupMounts() error {

  // mount the engine if it's a local directory
  if config.EngineDir() != "" {
    src := config.EngineDir()
    dst := fmt.Sprintf("%s%s/engine", provider.HostShareDir(), config.AppName())

    // first export the env on the workstation
    if err := envSetup.addShare(src, dst); err != nil {
      return err
    }

    // mount the env on the provider
    if err := envSetup.addMount(src, dst); err != nil {
      return err
    }
  }

  // mount the app src
  src := config.LocalDir()
  dst := fmt.Sprintf("%s%s/code", provider.HostShareDir(), config.AppName())

  // first export the env on the workstation
  if err := envSetup.addShare(src, dst); err != nil {
    return err
  }

  // then mount the env on the provider
  if err := envSetup.addMount(src, dst); err != nil {
    return err
  }

  return nil
}

// setupApp sets up the app plaftorm and data services
func (envSetup *processEnvSetup) setupApp() error {

  // setup the app
  if err := processor.Run("app_setup", envSetup.control); err != nil {
    return err
  }

  // clean up after any possible failures in a previous deploy
  if err := processor.Run("service_clean", envSetup.control); err != nil {
    return err
  }

  // setup the platform services
  return processor.Run("platform_setup", envSetup.control)
}

// addShare will add a filesystem env on the host machine
func (envSetup *processEnvSetup) addShare(src, dst string) error {

  // the mount type is configurable by the user
  mountType := config.Viper().GetString("mount-type")

  // todo: we should display a warning when using native about performance

  // since vm.mount is configurable, it's possible and even likely that a
  // machine may already have mounts configured. For each mount type we'll
  // need to check if an existing mount needs to be undone before continuing
  switch mountType {

  // check to see if netfs is currently configured. If it is then tear it down
  // and build the native env
  case "native":
    if netfs.Exists(src) {
      // netfs was used prior. So we need to tear it down.

      control := processor.ProcessControl{
        Env:          envSetup.control.Env,
        Verbose:      envSetup.control.Verbose,
        DisplayLevel: envSetup.control.DisplayLevel,
        Meta: map[string]string{
          "path": src,
        },
      }

      if err := processor.Run("env_netfs_remove", control); err != nil {
        return err
      }
    }

    // now we let the provider add it's native env
    if err := provider.AddShare(src, dst); err != nil {
      return err
    }

  // check to see if native envs are currently exported. If so,
  // tear down the native env and build the netfs env
  case "netfs":
    if provider.HasShare(src, dst) {
      // native was used prior. So we need to tear it down
      if err := provider.RemoveShare(src, dst); err != nil {
        return err
      }
    }

    control := processor.ProcessControl{
      Env:      envSetup.control.Env,
      Verbose:      envSetup.control.Verbose,
      DisplayLevel: envSetup.control.DisplayLevel,
      Meta: map[string]string{
        "path": src,
      },
    }

    if err := processor.Run("env_netfs_add", control); err != nil {
      return err
    }
  }

  return nil
}

// addMount will mount a env in the nanobox guest context
func (envSetup *processEnvSetup) addMount(src, dst string) error {

    // short-circuit if the mount already exists
    if provider.HasMount(dst) {
      return nil
    }

    // the mount type is configurable by the user
    mountType := config.Viper().GetString("mount-type")

    switch mountType {

    // build the native mount
    case "native":
      if err := provider.AddMount(src, dst); err != nil {
        return err
      }

    // build the netfs mount
    case "netfs":
      if err := netfs.Mount(src, dst); err != nil {
        return err
      }
    }

    return nil
}
