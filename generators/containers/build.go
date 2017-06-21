package containers

import (
	"fmt"

	"github.com/nanobox-io/golang-docker-client"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util/provider"
)

// BuildConfig generate the container configuration for the build container
func BuildConfig(image string) docker.ContainerConfig {
	env := config.EnvID()
	code := fmt.Sprintf("%s%s/code:/app", provider.HostShareDir(), env)
	engine := fmt.Sprintf("%s%s/engine:/share/engine", provider.HostShareDir(), env)

	if !provider.RequiresMount() {
		code = fmt.Sprintf("%s:/app", config.LocalDir())
		if config.EngineDir() != "" {
			engine = fmt.Sprintf("%s:/share/engine", config.EngineDir())
		}
	}

	cache := fmt.Sprintf("nanobox_%s_cache:/mnt/cache", env)
	configModel, _ := models.LoadConfig()
	if configModel.Cache == "shared" {
		cache = "nanobox_cache:/mtn/cache"
	}

	conf := docker.ContainerConfig{
		Name:    BuildName(),
		Image:   image,
		Network: "host",
		Binds: []string{
			code,
			engine,
			// fmt.Sprintf("%s%s/build:/mnt/build", provider.HostMntDir(), env),
			// fmt.Sprintf("%s%s/deploy:/mnt/deploy", provider.HostMntDir(), env),
			// fmt.Sprintf("%s%s/cache:/mnt/cache", provider.HostMntDir(), env),
			fmt.Sprintf("nanobox_%s_build:/mnt/build", env),
			fmt.Sprintf("nanobox_%s_deploy:/mnt/deploy", env),
			cache,
		},
		RestartPolicy: "no",
	}

	// Some CI's have an old kernel and require us to use the virtual network
	// this is only in effect for CI's because it automatically reserves an ip on our nanobox
	// virtual network and we could have IP conflicts
	configModel, _ := models.LoadConfig()
	if configModel.CIMode {
		conf.Network = "virt"
	}

	// set http[s]_proxy and no_proxy vars
	setProxyVars(&conf)

	if config.EngineDir() != "" {
		conf.Binds = append(conf.Binds, engine)
	}

	return conf
}

// BuildName returns the name of the build container
func BuildName() string {
	return fmt.Sprintf("nanobox_%s_build", config.EnvID())
}
