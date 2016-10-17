package containers

import (
	"fmt"

	"github.com/nanobox-io/golang-docker-client"
	"github.com/nanobox-io/nanobox-boxfile"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/dhcp"
	"github.com/nanobox-io/nanobox/util/provider"
)

// DevConfig generate the container configuration for the build container
func DevConfig(appModel *models.App) docker.ContainerConfig {
	boxfile := boxfile.New([]byte(appModel.DeployedBoxfile))

	image := boxfile.Node("build").StringValue("image")

	if image == "" {
		image = "nanobox/build"
	}

	config := docker.ContainerConfig{
		Name:    fmt.Sprintf("nanobox_%s", appModel.ID),
		Image:   image, // this will need to be configurable some time
		Network: "virt",
		IP:      reserveIP(),
		Binds: []string{
			fmt.Sprintf("%s%s/code:/app", provider.HostShareDir(), appModel.EnvID),
			fmt.Sprintf("%s%s/build:/data", provider.HostMntDir(), appModel.EnvID),
			fmt.Sprintf("%s%s/cache:/mnt/cache", provider.HostMntDir(), appModel.EnvID),
		},
	}

	// add lib_dirs into the container binds
	libDirs := boxfile.Node("code.build").StringSliceValue("lib_dirs")

	for _, libDir := range libDirs {
		// TODO: the cache source should come from the provider
		path := fmt.Sprintf("/mnt/sda1/%s/cache/lib_dirs/%s:/app/%s", appModel.EnvID, libDir, libDir)
		config.Binds = append(config.Binds, path)
	}

	return config
}

// reserveIP reserves a local IP for the build container
func reserveIP() string {
	ip, _ := dhcp.ReserveLocal()
	return ip.String()
}

// DevName returns the name of the build container
func DevName() string {
	return fmt.Sprintf("nanobox_%s_dev", config.EnvID())
}
