package processors

import (
	"strconv"

	"github.com/nanobox-io/nanobox/models"
)

func ConfigureSet(key, val string) error {
	config, _ := models.LoadConfig()

	switch key {
	case "provider":
		config.Provider = val
	case "mount-type", "mount_type":
		config.MountType = val
	case "netfs_mount_opts", "netfs-mount-opts", "mount_options", "mount-options":
		config.NetfsMountOpts = val
	case "cpus", "CPUs":
		config.CPUs, _ = strconv.Atoi(val)
	case "ram", "RAM":
		config.RAM, _ = strconv.Atoi(val)
	case "disk":
		config.Disk, _ = strconv.Atoi(val)
	case "external_network_space", "external-network-space":
		config.ExternalNetworkSpace = val
	case "docker_machine_network_space", "docker-machine-network-space":
		config.DockerMachineNetworkSpace = val
	case "native_network_space", "native-network-space":
		config.NativeNetworkSpace = val
	case "lock_port", "lock-port":
		config.LockPort, _ = strconv.Atoi(val)

	}

	return config.Save()
}
