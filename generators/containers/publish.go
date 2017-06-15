package containers

import (
	"fmt"

	"github.com/nanobox-io/golang-docker-client"

	"github.com/nanobox-io/nanobox/util/config"
	// "github.com/nanobox-io/nanobox/util/provider"
)

// PublishConfig generate the container configuration for the build container
func PublishConfig(image string) docker.ContainerConfig {
	env := config.EnvID()
	config := docker.ContainerConfig{
		Name:    PublishName(),
		Image:   image,
		Network: "virt",
		Binds: []string{
			// fmt.Sprintf("%s%s/app:/mnt/app", provider.HostMntDir(), env),
			// fmt.Sprintf("%s%s/cache:/mnt/cache", provider.HostMntDir(), env),
			// fmt.Sprintf("%s%s/deploy:/mnt/deploy", provider.HostMntDir(), env),
			fmt.Sprintf("nanobox_%s_app:/mnt/app", env),
			fmt.Sprintf("nanobox_%s_cache:/mnt/cache", env),
			fmt.Sprintf("nanobox_%s_deploy:/mnt/deploy", env),
		},
		RestartPolicy: "no",
	}

	// set http[s]_proxy and no_proxy vars
	setProxyVars(&config)

	return config
}

// PublishName returns the name of the build container
func PublishName() string {
	return fmt.Sprintf("nanobox_%s_publish", config.EnvID())
}
