package containers

import (
	"fmt"
	"os"

	"github.com/nanobox-io/golang-docker-client"
)

func setProxyVars(config *docker.ContainerConfig) {
	// set the proxy variables
	httpProxyEvar := os.Getenv("HTTP_PROXY")
	if httpProxyEvar != "" {
		config.Env = append(config.Env, fmt.Sprintf("HTTP_PROXY=%s", httpProxyEvar))
	}
	httpsProxyEvar := os.Getenv("HTTPS_PROXY")
	if httpsProxyEvar != "" {
		config.Env = append(config.Env, fmt.Sprintf("HTTPS_PROXY=%s", httpsProxyEvar))
	}
	noProxyEvar := os.Getenv("NO_PROXY")
	if noProxyEvar != "" {
		config.Env = append(config.Env, fmt.Sprintf("NO_PROXY=%s", noProxyEvar))
	}
	httpProxyEvar2 := os.Getenv("http_proxy")
	if httpProxyEvar2 != "" {
		config.Env = append(config.Env, fmt.Sprintf("http_proxy=%s", httpProxyEvar2))
	}
	httpsProxyEvar2 := os.Getenv("https_proxy")
	if httpsProxyEvar2 != "" {
		config.Env = append(config.Env, fmt.Sprintf("https_proxy=%s", httpsProxyEvar2))
	}
	noProxyEvar2 := os.Getenv("no_proxy")
	if noProxyEvar2 != "" {
		config.Env = append(config.Env, fmt.Sprintf("no_proxy=%s", noProxyEvar2))
	}
}
