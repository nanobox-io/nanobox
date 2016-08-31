package containers

import (
	"fmt"

	"github.com/nanobox-io/golang-docker-client"

	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/provider"
)

// BuildConfig generate the container configuration for the build container
func BuildConfig(image, ip string) docker.ContainerConfig {
	env := config.EnvID()
	config := docker.ContainerConfig{
		Name:    BuildName(),
		Image:   image,
		Network: "virt",
		IP:      ip,
		Binds: []string{
			fmt.Sprintf("%s%s/code:/share/code", provider.HostShareDir(), env),
			fmt.Sprintf("%s%s/engine:/share/engine", provider.HostShareDir(), env),
			fmt.Sprintf("%s%s/build:/mnt/build", provider.HostMntDir(), env),
			fmt.Sprintf("%s%s/deploy:/mnt/deploy", provider.HostMntDir(), env),
			fmt.Sprintf("%s%s/app:/mnt/app", provider.HostMntDir(), env),
			fmt.Sprintf("%s%s/cache:/mnt/cache", provider.HostMntDir(), env),
		},
	}

	return config
}

// BuildName returns the name of the build container
func BuildName() string {
	return fmt.Sprintf("nanobox_%s_build", config.EnvID())
}
