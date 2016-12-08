package containers

import (
	"github.com/nanobox-io/golang-docker-client"

	"github.com/nanobox-io/nanobox/util/dhcp"
)

// BridgeConfig generates the container configuration for a component container
func BridgeConfig() docker.ContainerConfig {
	return docker.ContainerConfig{
		Name:          BridgeName(),
		Image:         "nanobox/bridge",
		Network:       "virt",
		IP:            reserveIP(),
		RestartPolicy: "always",
		Ports:         []string{"1194:1194/udp"},
	}
}

// BridgeName returns the name of the component container
func BridgeName() string {
	return "nanobox_bridge"
}


// reserveIP reserves a local IP for the build container
func reserveIP() string {
	ip, _ := dhcp.ReserveLocal()
	return ip.String()
}