package docker

import (
	dockClient "github.com/docker/engine-api/client"
)

var client *dockClient.Client

func init() {
	Initialize("env")
}

func Initialize(host string) error {
	var err error
	if host == "env" {
		client, err = dockClient.NewEnvClient()
	} else {
		client, err = dockClient.NewClient(host, "", nil, nil)
	}
	if err != nil {
		client = nil
		return err
	}
	// this wasnt being used... i dont think
	// networks, err := client.NetworkList(context.Background(), dockType.NetworkListOptions{})
	// if err == nil {
	// 	for _, network := range networks {
	// 		if network.Name == "nanobox" {
	// 			viper.Set("nanobox-network", network.ID)
	// 			break
	// 		}
	// 	}
	// }
	return nil
}
