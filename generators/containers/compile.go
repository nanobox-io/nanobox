package containers

import (
	"fmt"

	"github.com/nanobox-io/golang-docker-client"

	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/provider"
)

// CompileConfig generate the container configuration for the build container
func CompileConfig(image string) docker.ContainerConfig {
	env := config.EnvID()
	config := docker.ContainerConfig{
		Name:    CompileName(),
		Image:   image,
		Network: "host",
		Binds: []string{
			fmt.Sprintf("%s%s/code:/share/code", provider.HostShareDir(), env),
			fmt.Sprintf("%s%s/engine:/share/engine", provider.HostShareDir(), env),
			fmt.Sprintf("%s%s/build:/data", provider.HostMntDir(), env),
			fmt.Sprintf("%s%s/app:/mnt/app", provider.HostMntDir(), env),
			fmt.Sprintf("%s%s/cache:/mnt/cache", provider.HostMntDir(), env),
		},
		RestartPolicy: "no",
	}

	return config
}

// CompileName returns the name of the build container
func CompileName() string {
	return fmt.Sprintf("nanobox_%s_compile", config.EnvID())
}
