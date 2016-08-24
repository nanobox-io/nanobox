package code

import (
  "fmt"
  
  "github.com/jcelliott/lumber"
  "github.com/nanobox-io/golang-docker-client"
  
  "github.com/nanobox-io/nanobox/commands/registry"
  "github.com/nanobox-io/nanobox/generators/containers"
  "github.com/nanobox-io/nanobox/models"
  "github.com/nanobox-io/nanobox/processors/env"
  "github.com/nanobox-io/nanobox/util/boxfile"
  "github.com/nanobox-io/nanobox/util/display"
  "github.com/nanobox-io/nanobox/util/dhcp"
)

// Build builds the codebase that can then be deployed
func Build(env *models.Env) error {
  // extract the build image from the boxfile
  buildImage := boxfile.BuildImage()
  
  // if a build container was leftover from a previous deploy, let's remove it
  docker.ContainerRemove(containers.BuildContainerName())
  
  // generate a docker percent display
  dockerPercent := &display.DockerPercentDisplay{
    Output: display.NewStreamer("info"), 
    Prefix: buildImage,
  }
  
  // pull the build image
  if _, err := docker.ImagePull(buildImage, dockerPercent); err != nil {
    lumber.Error("code:Build:docker.ImagePull(%s, nil): %s", buildImage, err.Error())
    return fmt.Errorf("failed to pull docker image (%s): %s", buildImage, err.Error())
  }
  
  // reserve an IP for the build container
  ip, err := dhcp.ReserveLocal()
  if err != nil {
    lumber.Error("code:Build:dhcp.ReserveLocal(): %s", err.Error())
    return fmt.Errorf("failed to reserve an ip for the build container: %s", err.Error())
  }
  
  // ensure we release the IP when we're done
  defer dhcp.ReturnIP(ip)
  
  // start the container
  config := containers.BuildConfig(buildImage, ip.String())
  container, err := docker.CreateContainer(config)
  if err != nil {
    lumber.Error("code:Build:docker.CreateContainer(%+v): %s", config, err.Error())
    return fmt.Errorf("failed to start docker container: %s", err.Error())
  }
  
  // ensure we stop the container when we're done
  defer docker.ContainerRemove(containers.BuildContainerName())
  
  // run the user hook
  if _, err := RunUserHook(container.ID); err != nil {
    err = fmt.Errorf("failed to run user hook: %s", err.Error())
    return runBuildDebugSession(container.ID, err)
  }
  
  // run the configure hook
  if _, err := RunConfigureHook(container.ID); err != nil {
    err = fmt.Errorf("failed to run configure hook: %s", err.Error())
    return runBuildDebugSession(container.ID, err)
  }
  
  // run the fetch hook
  if _, err := RunFetchHook(container.ID); err != nil {
    err = fmt.Errorf("failed to run fetch hook: %s", err.Error())
    return runBuildDebugSession(container.ID, err)
  }
  
  // run the setup hook
  if _, err := RunSetupHook(container.ID); err != nil {
    err = fmt.Errorf("failed to run setup hook: %s", err.Error())
    return runBuildDebugSession(container.ID, err)
  }
  
  // run the boxfile hook
  boxOutput, err := RunBoxfileHook(container.ID)
  if err != nil {
    err = fmt.Errorf("failed to run boxfile hook: %s", err.Error())
    return runBuildDebugSession(container.ID, err)
  }
  
  // persist the boxfile output to the env model
  env.BuiltBoxfile = boxOutput
  if err := env.Save(); err != nil {
    lumber.Error("code:Build:models:Env:Save(): %s", err.Error())
    return fmt.Errorf("failed to persist build boxfile to db: %s", err.Error())
  }
  
  // run the prepare hook
  if _, err := RunPrepareHook(container.ID); err != nil {
    err = fmt.Errorf("failed to run prepare hook: %s", err.Error())
    return runBuildDebugSession(container.ID, err)
  }
  
  if !registry.GetBool("skip-compile") {
    // run the compile hook
    if _, err := RunCompileHook(container.ID); err != nil {
      err = fmt.Errorf("failed to run compile hook: %s", err.Error())
      return runBuildDebugSession(container.ID, err)
    }
    
    // run the pack-app hook
    if _, err := RunPackAppHook(container.ID); err != nil {
      err = fmt.Errorf("failed to run pack-app hook: %s", err.Error())
      return runBuildDebugSession(container.ID, err)
    }
  }
  
  // run the pack-build hook
  if _, err := RunPackBuildHook(container.ID); err != nil {
    err = fmt.Errorf("failed to run pack-build hook: %s", err.Error())
    return runBuildDebugSession(container.ID, err)
  }
  
  if !registry.GetBool("skip-compile") {
    // run the clean hook
    if _, err := RunCleanHook(container.ID); err != nil {
      err = fmt.Errorf("failed to run clean hook: %s", err.Error())
      return runBuildDebugSession(container.ID, err)
    }
    
    // run the pack-deploy hook
    if _, err := RunPackDeployHook(container.ID); err != nil {
      err = fmt.Errorf("failed to run pack-deploy hook: %s", err.Error())
      return runBuildDebugSession(container.ID, err)
    }
  }
  
  return nil
}


// runDebugSession drops the user in the build container to debug
func runBuildDebugSession(container string, err error) error {
	if registry.GetBool("debug") {
		component := models.Component{
			ID: container,
		}
		envConsole := env.Console{
			Component: component,
		}
		err := envConsole.Run()
		if err != nil {
			return fmt.Errorf("failed to establish a debug session: %s", err.Error())
		}
	}

  return err
}
