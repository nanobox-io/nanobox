package docker

import (
	"golang.org/x/net/context"
	"strings"
	"time"

	dockType "github.com/docker/engine-api/types"
	dockContainer "github.com/docker/engine-api/types/container"
	dockNetwork "github.com/docker/engine-api/types/network"
	"github.com/docker/engine-api/types/strslice"
)

type ContainerConfig struct {
	ID         string            `json:"id"`
	Network    string            `json:"network"`
	Name       string            `json:"name"`
	Labels     map[string]string `json:"labels"`
	Hostname   string            `json:"hostname"`
	Domainname string            `json:"domainname"`
	Cmd        []string          `json:"cmd"`
	Image      string            `json:"image_slug"`
	IP         string            `json:"ip"`
	Binds      []string          `json:"binds"`
	Memory     int64             `json:"memory"`
	MemorySwap int64             `json:"memory_swap"`
	Status     string            `json:"status"`
	CPUShares  int64             `json:"cpu_shares"`
}

// create a container from the user specification
func CreateContainer(conf ContainerConfig) (dockType.ContainerJSON, error) {
	// if len(conf.Cmd) == 0 {
	// 	conf.Cmd = []string{"/bin/sleep", "3650d"}
	// }

	config := &dockContainer.Config{
		Hostname:        conf.Hostname,
		Domainname:      conf.Domainname,
		Cmd:             conf.Cmd,
		Labels:          conf.Labels,
		NetworkDisabled: false,
		Image:           conf.Image,
	}

	hostConfig := &dockContainer.HostConfig{
		Privileged: true,
		Binds:      conf.Binds,
		// NetworkMode:   "host",
		CapAdd:        strslice.StrSlice([]string{"NET_ADMIN"}),
		NetworkMode:   "bridge",
		RestartPolicy: dockContainer.RestartPolicy{Name: "unless-stopped"},
		Resources: dockContainer.Resources{
			Memory:     conf.Memory,
			MemorySwap: conf.MemorySwap,
			CPUShares:  conf.CPUShares,
		},
	}

	netConfig := &dockNetwork.NetworkingConfig{}

	if conf.Network == "host" {
		// you cant set the hostname of the host
		hostConfig.NetworkMode = "host"
		config.Hostname = ""
	}

	if conf.Network == "virt" || conf.IP != "" {
		hostConfig.NetworkMode = "nanobox"
		netConfig.EndpointsConfig = map[string]*dockNetwork.EndpointSettings{
			"nanobox": &dockNetwork.EndpointSettings{
				IPAMConfig: &dockNetwork.EndpointIPAMConfig{IPv4Address: conf.IP},
			},
		}
	}

	return createContainer(config, hostConfig, netConfig, conf.Name)
}

// createContainer
func createContainer(config *dockContainer.Config, hostConfig *dockContainer.HostConfig, networkingConfig *dockNetwork.NetworkingConfig, containerName string) (dockType.ContainerJSON, error) {
	// create container
	container, err := client.ContainerCreate(context.Background(), config, hostConfig, networkingConfig, containerName)
	if err != nil {
		return dockType.ContainerJSON{}, err
	}

	if err := ContainerStart(container.ID); err != nil {
		return dockType.ContainerJSON{}, err
	}

	return ContainerInspect(container.ID)
}

// Start a container. If the container is already running this will error.
func ContainerStart(id string) error {
	return client.ContainerStart(context.Background(), id, dockType.ContainerStartOptions{})
}

// give them 5 seconds.. todo maybe make it adjustable
func ContainerStop(id string) error {
	timeout := 5 * time.Second
	return client.ContainerStop(context.Background(), id, &timeout)
}

// ContainerRemove
func ContainerRemove(id string) error {
	_, err := client.ContainerInspect(context.Background(), id)
	if err != nil {
		return err
	}

	timeout := 0 * time.Second
	if err := client.ContainerStop(context.Background(), id, &timeout); err != nil {
		// return err
	}

	return client.ContainerRemove(context.Background(), id, dockType.ContainerRemoveOptions{RemoveVolumes: true, Force: true})
}

// ContainerInspect
func ContainerInspect(id string) (dockType.ContainerJSON, error) {
	return client.ContainerInspect(context.Background(), id)
}

// GetContainer
func GetContainer(name string) (dockType.ContainerJSON, error) {
	return client.ContainerInspect(context.Background(), name)
}

// ContainerList
func ContainerList() ([]dockType.Container, error) {
	return client.ContainerList(context.Background(), dockType.ContainerListOptions{All: true, Size: false})
}

func ContainerJSONtoConfig(cj dockType.ContainerJSON) ContainerConfig {

	return ContainerConfig{
		ID:         cj.ID,
		Network:    string(cj.HostConfig.NetworkMode),
		Name:       strings.Replace(cj.Name, "/", "", 1),
		Hostname:   cj.Config.Hostname,
		Domainname: cj.Config.Domainname,
		Labels:     cj.Config.Labels,
		Cmd:        []string(cj.Config.Cmd),
		Image:      cj.Config.Image,
		IP:         GetIP(cj),
		Memory:     cj.HostConfig.Resources.Memory,
		MemorySwap: cj.HostConfig.Resources.MemorySwap,
		Status:     cj.State.Status,
		CPUShares:  cj.HostConfig.Resources.CPUShares,
	}
}

func ContainerSliceToConfigSlice(cs []dockType.Container) []ContainerConfig {
	cslice := []ContainerConfig{}
	for _, c := range cs {

		name := ""
		if len(c.Names) >= 1 {
			name = c.Names[0]
		}
		cslice = append(cslice, ContainerConfig{
			ID:      c.ID,
			Network: c.HostConfig.NetworkMode,
			Name:    strings.Replace(name, "/", "", 1),
			Image:   c.Image,
			Status:  c.Status,
			Labels:  c.Labels,
			IP:      GetIP(cs),
		})
	}
	return cslice
}

func GetIP(i interface{}) (ip string) {
	switch i.(type) {
	default:

	case dockType.Container:
		c := i.(dockType.Container)
		if c.NetworkSettings != nil {
			for _, val := range c.NetworkSettings.Networks {
				if val.IPAddress != "" {
					ip = val.IPAddress
					break
				}
			}
		}
	case dockType.ContainerJSON:
		c := i.(dockType.ContainerJSON)
		ip = c.NetworkSettings.IPAddress
		if ip == "" {
			for _, val := range c.NetworkSettings.Networks {
				if val.IPAddress != "" {
					ip = val.IPAddress
					break
				}
			}
		}
	}
	return
}
