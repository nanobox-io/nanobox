package dev

import (
  "fmt"
  "net"
  "os"

  "github.com/jcelliott/lumber"
  "github.com/nanobox-io/golang-docker-client"
  "github.com/nanobox-io/nanobox-boxfile"

  container_generator "github.com/nanobox-io/nanobox/generators/containers"
  build_generator     "github.com/nanobox-io/nanobox/generators/hooks/build"
  
  "github.com/nanobox-io/nanobox/models"
  "github.com/nanobox-io/nanobox/processors/env"
  "github.com/nanobox-io/nanobox/util/counter"
  "github.com/nanobox-io/nanobox/util/dhcp"
  "github.com/nanobox-io/nanobox/util/display"
  "github.com/nanobox-io/nanobox/util/hookit"
  "github.com/nanobox-io/nanobox/util/locker"
  "github.com/nanobox-io/nanobox/util/provider"
  // "github.com/nanobox-io/nanobox/util/watch"
)

// Start a dev container 
func Console(envModel *models.Env, appModel *models.App, devRun bool) error {
  
  // ensure the environment is setup
  if err := env.Setup(envModel); err != nil {
    return fmt.Errorf("failed to setup environment: %s", err.Error())
  }
  
  // whatever happens next, ensure we teardown this container
  defer teardown(appModel)
  
  // setup the dev container
  if err := setup(appModel); err != nil {
    return fmt.Errorf("failed to setup dev container: %s", err.Error())
  }
  
  // start a watcher to watch for changes and inform the vm
  // go watch.Watch(container_generator.DevName(), envModel.Directory)
  
  // if run then start the run commands
  if devRun {
    return Run(appModel)
  }
  
  // print the MOTD before dropping into the container
  if err := printMOTD(appModel); err != nil {
    return fmt.Errorf("failed to print MOTD: %s", err.Error())
  }

  // console into the newly created container
  if err := runConsole(appModel); err != nil {
    return fmt.Errorf("failed to console into dev container: %s", err.Error())
  }

  return nil
}

// sets up the dev container and network stack
func setup(appModel *models.App) error {

	// establish a local lock to ensure we're the only ones bringing up the
	// dev container. Also, we need to ensure the lock is released even in we error
	locker.LocalLock()
	defer locker.LocalUnlock()
  
	// let anyone else know we're using the dev container
	counter.Increment(appModel.ID)

  // we don't need to setup if dev is already running
  if isDevRunning() {
    return nil
  }

  display.OpenContext("Building dev environment")
  defer display.CloseContext()

  // generate a container config
  config := container_generator.DevConfig(appModel)

	//
	if err := downloadImage(config.Image); err != nil {
    return err
	}

  display.StartTask("Starting docker container")
	container, err := docker.CreateContainer(config)
	if err != nil {
    display.ErrorTask()
		return fmt.Errorf("failed to create docker container: %s", err.Error())
	}
  display.StopTask()

	//
	if err := attachNetwork(appModel, config.IP); err != nil {
		return fmt.Errorf("unable to attach container to network: %s", err.Error())
	}

	lumber.Prefix("dev:Console")
	defer lumber.Prefix("")

  display.StartTask("Configuring")
	userPayload := build_generator.UserPayload()
	if _, err := hookit.Exec(container.ID, "user", userPayload, "debug"); err != nil {
    display.ErrorTask()
		return fmt.Errorf("failed to run the user hook: %s", err.Error())
	}

	if _, err := hookit.Exec(container.ID, "dev", build_generator.DevPayload(appModel), "debug"); err != nil {
    display.ErrorTask()
		return fmt.Errorf("failed to run the dev hook: %s", err.Error())
	}
  display.StopTask()

	return nil
}

func teardown(appModel *models.App) error {
  locker.LocalLock()
  defer locker.LocalUnlock()

  counter.Decrement(appModel.ID)
  
  if !devIsUnused(appModel.ID) {
    return nil
  }
  
  // grab the container info
	container, err := docker.GetContainer(container_generator.DevName())
	if err != nil {
		// if we cant get the container it may have been removed by someone else
		// just return here
		return nil
	}
  
  // remove the container
  if err := docker.ContainerRemove(container.ID); err != nil {
		return fmt.Errorf("failed to remove dev container: %s", err.Error())
	}
  
  // extract the container IP
  ip := docker.GetIP(container)
  
  // detach dev container from the network
  if err := detachNetwork(appModel, ip); err != nil {
    return fmt.Errorf("failed to detach dev container from network: %s", err.Error())
  }
  
  // return the container IP back to the IP pool
  if err := dhcp.ReturnIP(net.ParseIP(ip)); err != nil {
    lumber.Error("An error occurred durring dev console teadown:%s", err.Error())
    return fmt.Errorf("failed to return unused IP back to pool: %s", err.Error())
  }
  
  return nil
}

// attachNetwork attaches the container to the host network
func attachNetwork(appModel *models.App, containerIP string) error {
  display.StartTask("Attaching network")
	defer display.StopTask()

	// fetch the devIP
	devIP := appModel.GlobalIPs["env"]

	//
	if err := provider.AddIP(devIP); err != nil {
    lumber.Error("dev:attachNetwork:provider.AddIP(%s):%s", devIP, err.Error())
		return fmt.Errorf("failed to add IP to the provider: %s", err.Error())
	}

	//
	if err := provider.AddNat(devIP, containerIP); err != nil {
    lumber.Error("dev:attachNetwork:provider.AddNat(%s, %s):%s", devIP, containerIP, err.Error())
		return fmt.Errorf("failed to add NAT from container: %s", err.Error())
	}

	return nil
}

// detachNetwork detaches the container from the host network
func detachNetwork(appModel *models.App, containerIP string) error {

	// fetch the devIP
	devIP := appModel.GlobalIPs["env"]

	// remove nat
	if err := provider.RemoveNat(devIP, containerIP); err != nil {
    lumber.Error("dev:detachNetwork:provider.RemoveNat(%s, %s):%s", devIP, containerIP, err.Error())
		return fmt.Errorf("failed to remove NAT from container: %s", err.Error())
	}

	// remove the IP from the provider
	if err := provider.RemoveIP(devIP); err != nil {
    lumber.Error("dev:detachNetwork:provider.RemoveIP(%s):%s", devIP, err.Error())
		return fmt.Errorf("failed to remove the IP from the provider: %s", err.Error())
	}

	return nil
}

// downloadImage downloads the dev docker image
func downloadImage(image string) error {

  display.StartTask("Pulling %s image", image)
  defer display.StopTask()

  // generate a docker percent display
  dockerPercent := &display.DockerPercentDisplay{
    Output: display.NewStreamer("info"),
    Prefix: image,
  }

	if _, err := docker.ImagePull(image, dockerPercent); err != nil {
    display.ErrorTask()
    lumber.Error("dev:Setup:downloadImage.ImagePull(%s, nil): %s", image, err.Error())
    return fmt.Errorf("failed to pull docker image (%s): %s", image, err.Error())
	}

	return nil
}

// printMOTD prints the motd with information for the user to connect
func printMOTD(appModel *models.App) error {

	// fetch the dev ip
	devIP := appModel.GlobalIPs["env"]

	os.Stderr.WriteString(fmt.Sprintf(`
                                   **
                                ********
                             ***************
                          *********************
                            *****************
                          ::    *********    ::
                             ::    ***    ::
                           ++   :::   :::   ++
                              ++   :::   ++
                                 ++   ++
                                    +
                    _  _ ____ _  _ ____ ___  ____ _  _
                    |\ | |__| |\ | |  | |__) |  |  \/
                    | \| |  | | \| |__| |__) |__| _/\_

--------------------------------------------------------------------------------
+ You are in a virtual machine (vm)
+ Your local source code has been mounted into the vm
+ Changes to your code in either the vm or workstation will be mirrored
+ If you run a server, access it at >> %s
--------------------------------------------------------------------------------

`, devIP))

	return nil
}

// runConsole will establish a console within the dev container
func runConsole(appModel *models.App) error {

	// create a dummy component using the appname
	component := &models.Component{
		ID: "nanobox_" + appModel.ID,
	}
  
	consoleConfig := env.ConsoleConfig{
		Cwd: cwd(appModel),
	}

	return env.Console(component, consoleConfig)
}

// cwd sets the cwd from the boxfile or provides a sensible default
func cwd(appModel *models.App) string {
	boxfile := boxfile.New([]byte(appModel.DeployedBoxfile))
	// parse the boxfile data

	if boxfile.Node("dev").StringValue("cwd") != "" {
		return boxfile.Node("dev").StringValue("cwd")
	}

	return "/app"
}

// devIsUnused returns true if the dev isn't being used by any other session
func devIsUnused(ID string) bool {
	count, err := counter.Get(ID)
	return count == 0 && err == nil
}

// isDevRunning returns true if a service is already running
func isDevRunning() bool {

	container, err := docker.GetContainer(container_generator.DevName())

	// if the container doesn't exist then just return false
	return err == nil && container.State.Status == "running"
}
