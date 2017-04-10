package bridge

import (
	"os"
	// "runtime"
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"time"

	"github.com/jcelliott/lumber"
	"github.com/nanobox-io/golang-docker-client"

	container_generator "github.com/nanobox-io/nanobox/generators/containers"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/display"
	// "github.com/nanobox-io/nanobox/util/fileutil"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/hookit"
	"github.com/nanobox-io/nanobox/util/locker"
)

var keys map[string]string

func Setup() error {

	display.OpenContext("Building bridge")
	defer display.CloseContext()

	// if the container exists and openvpn is running
	if Connected() {
		return nil
	}

	// create a component
	if err := setupContainer(); err != nil {
		return err
	}

	// configure bridge client
	if err := configureBridge(); err != nil {
		return err
	}

	// start bridge client
	if err := Start(); err != nil {
		return err
	}

	return nil
}

// sets up the dev container and network stack
func setupContainer() error {

	// we don't need to setup if bridge is already running
	if isBridgeExists() {
		return nil
	}

	// establish a local lock to ensure we're the only ones bringing up the
	// dev container. Also, we need to ensure the lock is released even in we error
	locker.LocalLock()
	defer locker.LocalUnlock()

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
		return util.ErrorAppend(err, "failed to create docker container")
	}
	display.StopTask()

	display.StartTask("Setting up container")
	// run the configure hook
	if _, err := hookit.DebugExec(container.ID, "configure", "{\"platform\":\"local\",\"config\":{}}", "info"); err != nil {
		return util.ErrorAppend(err, "failed to run configure hook")
	}

	// run the start hook
	if _, err := hookit.DebugExec(container.ID, "start", "{}", "info"); err != nil {
		return util.ErrorAppend(err, "failed to run start hook")
	}

	// run the start hook
	output, err := hookit.DebugExec(container.ID, "keys", "{}", "info")
	if err != nil {
		return util.ErrorAppend(err, "failed to run start hook")
	}

	if err := json.Unmarshal([]byte(output), &keys); err != nil {
		return util.ErrorAppend(err, "failed to decode the keys")
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
		return util.ErrorAppend(err, "failed to pull docker image (%s)", image)
	}

	return nil
}

// ca.crt
// client.key
// client.crt
func configureBridge() error {
	display.StartTask("Configuring")
	defer display.StopTask()

	// make the openvpn folder
	vpnDir := filepath.ToSlash(filepath.Join(config.EtcDir(), "openvpn"))

	if err := os.MkdirAll(vpnDir, 0755); err != nil {
		lumber.Fatal("[bridge] os.Mkdir() failed", err.Error())
	}

	// create config file
	if err := ioutil.WriteFile(ConfigFile(), []byte(BridgeConfig()), 0644); err != nil {
		return err
	}

	// make sure to not overwrite the keys if we didnt create the container on this run
	if keys["ca.crt"] == "" {
		return nil
	}

	// create ca.crt
	if err := ioutil.WriteFile(CaCrt(), []byte(keys["ca.crt"]), 0644); err != nil {
		return err
	}

	// create client.key
	if err := ioutil.WriteFile(ClientKey(), []byte(keys["client.key"]), 0644); err != nil {
		return err
	}
	// create client.crt
	if err := ioutil.WriteFile(ClientCrt(), []byte(keys["client.crt"]), 0644); err != nil {
		return err
	}

	return nil
}
