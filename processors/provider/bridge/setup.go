package bridge

import (
	"fmt"
	"os"
	"runtime"
	"time"
	"encoding/json"

	"github.com/jcelliott/lumber"
	"github.com/nanobox-io/golang-docker-client"

	container_generator "github.com/nanobox-io/nanobox/generators/containers"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/display"
	"github.com/nanobox-io/nanobox/util/fileutil"
	"github.com/nanobox-io/nanobox/util/hookit"
	"github.com/nanobox-io/nanobox/util/locker"
)

// ca.crt
// client.key
// client.crt
var keys map[string]string

func Setup() error {
	// create a component
	if err := setupContainer(); err != nil {
		return err
	}

	// download bridge client
	if err := downloadBridgeClient(); err != nil {
		return err
	}

	// configure bridge client
	if err := configureBridge(); err != nil {
		return err
	}

	// start bridge client
	if err := startBridge(); err != nil {
		return err
	}

	return nil
}

// sets up the dev container and network stack
func setupContainer() error {

	// establish a local lock to ensure we're the only ones bringing up the
	// dev container. Also, we need to ensure the lock is released even in we error
	locker.LocalLock()
	defer locker.LocalUnlock()

	// we don't need to setup if dev is already running
	if isBridgeExists() {
		return nil
	}

	display.OpenContext("Building bridge")
	defer display.CloseContext()

	// generate a container config
	config := container_generator.BridgeConfig()

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

	display.StartTask("Configuring")
	// run the configure hook
	if _, err := hookit.DebugExec(container.ID, "configure", "{\"platform\":\"local\",\"config\":{}}", "info"); err != nil {
		return fmt.Errorf("failed to run configure hook: %s", err.Error())
	}

	// run the start hook
	if _, err := hookit.DebugExec(container.ID, "start", "{}", "info"); err != nil {
		return fmt.Errorf("failed to run start hook: %s", err.Error())
	}

	// run the start hook
	output, err := hookit.DebugExec(container.ID, "keys", "{}", "info")
	if err != nil {
		return fmt.Errorf("failed to run start hook: %s, %s", output, err.Error())
	}
	
	fmt.Println("keys output", output)
	if err := json.Unmarshal([]byte(output), &keys); err != nil {
		return fmt.Errorf("failed to decode the keys: %s %s", output, err.Error())
	}
	display.StopTask()

	return nil
}

// isBridgeExists returns true if a service is already running
func isBridgeExists() bool {

	_, err := docker.GetContainer(container_generator.BridgeName())

	// if the container doesn't exist then just return false
	return err == nil
}

// downloadImage downloads the dev docker image
func downloadImage(image string) error {

	if docker.ImageExists(image) {
		return nil
	}

	display.StartTask("Pulling %s image", image)
	defer display.StopTask()

	// generate a docker percent display
	dockerPercent := &display.DockerPercentDisplay{
		Output: display.NewStreamer("info"),
		// Prefix: image,
	}

	imagePull := func() error {
		_, err := docker.ImagePull(image, dockerPercent)
		return err
	}
	if err := util.Retry(imagePull, 5, time.Second); err != nil {
		display.ErrorTask()
		lumber.Error("dev:Setup:downloadImage.ImagePull(%s, nil): %s", image, err.Error())
		return fmt.Errorf("failed to pull docker image (%s): %s", image, err.Error())
	}

	return nil
}

func downloadBridgeClient() error {
	// TODO: remove once we have a client we can download
	return nil
	// short-circuit if we're already installed
	if fileutil.Exists(bridgeClient) {
		return nil
	}

	display.StartTask("Downloading bridge client")
	defer display.StopTask()

	// download the executable
	if err := fileutil.Download(bridgeURL, bridgeClient); err != nil {
		display.ErrorTask()
		return fmt.Errorf("failed to download docker-machine: %s", err.Error())
	}

	// make it executable (unless it's windows)
	if runtime.GOOS != "windows" {
		// make new CLI executable
		if err := os.Chmod(bridgeClient, 0755); err != nil {
			display.ErrorTask()
			return fmt.Errorf("failed to set permissions: ", err.Error())
		}
	}

	return nil
}

func configureBridge() error {

	return nil
}

func startBridge() error {

	return nil
}
