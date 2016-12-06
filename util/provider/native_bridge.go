package provider

import (
	"runtime"
)

func (native Native) createBridge() error {
	// linux doesnt need to bridge networks
	if runtime.GOOS == "linux" {
		return nil
	}

	// create docker container

	// setup vpn client

	return nil
}