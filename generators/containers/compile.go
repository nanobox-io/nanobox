package containers

import (
	"fmt"

	"github.com/nanobox-io/golang-docker-client"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/provider"
)

// CompileConfig generate the container configuration for the build container
func CompileConfig(image string) docker.ContainerConfig {
	env := config.EnvID()
	code := fmt.Sprintf("%s%s/code:/share/code", provider.HostShareDir(), env)
	engine := fmt.Sprintf("%s%s/engine:/share/engine", provider.HostShareDir(), env)

	if !provider.RequiresMount() {
		code = fmt.Sprintf("%s:/share/code", config.LocalDir())
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
		Name:    CompileName(),
		Image:   image,
		Network: "host",
		Binds: []string{
			code,
			engine,
			// fmt.Sprintf("%s%s/build:/data", provider.HostMntDir(), env),
			// fmt.Sprintf("%s%s/app:/mnt/app", provider.HostMntDir(), env),
			// fmt.Sprintf("%s%s/cache:/mnt/cache", provider.HostMntDir(), env),
			fmt.Sprintf("nanobox_%s_build:/data", env),
			fmt.Sprintf("nanobox_%s_app:/mnt/app", env),
			cache,
		},
		RestartPolicy: "no",
	}

	if config.EngineDir() != "" {
		conf.Binds = append(conf.Binds, engine)
	}

	return conf
}

// CompileName returns the name of the build container
func CompileName() string {
	return fmt.Sprintf("nanobox_%s_compile", config.EnvID())
}
