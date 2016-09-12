package nanoagent

import (
	"fmt"
	"github.com/nanobox-io/golang-ssh"
	"strconv"
	"strings"
)

func SSH(key, location string) error {

	// create the ssh client
	nanPass := ssh.Auth{Passwords: []string{key}}
	locationParts := strings.Split(location, ":")
	if len(locationParts) != 2 {
		return fmt.Errorf("location is not formatted properly (%s)", location)
	}

	// parse port
	port, err := strconv.Atoi(locationParts[1])
	if err != nil {
		return fmt.Errorf("unable to convert port (%s)", locationParts[1])
	}

	// establish connection
	client, err := ssh.NewNativeClient(key, locationParts[0], "SSH-2.0-nanobox", port, &nanPass)
	if err != nil {
		return fmt.Errorf("Failed to create new client - %s", err)
	}

	// establish the ssh client connection and shell
	err = client.Shell()
	if err != nil && err.Error() != "exit status 255" {
		return fmt.Errorf("Failed to request shell - %s", err)
	}

	return nil
}
