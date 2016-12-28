package docker

import (
	dockType "github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/filters"
	"golang.org/x/net/context"


)

// create a new volume
func VolumeCreate(name string) (dockType.Volume, error) {
	vol := dockType.VolumeCreateRequest{
		Name:       name,
	}

	return client.VolumeCreate(context.Background(), vol)
}

// list the volumes we have
func VolumeList() ([]*dockType.Volume, error) {
	volumeList, err := client.VolumeList(context.Background(), filters.Args{})
	return volumeList.Volumes, err
}

// check to see if a volume exists
func VolumeExists(name string) bool {
	volumes, err := VolumeList()
	if err != nil {
		return false
	}

	for _, volume := range volumes {
		if volume.Name == name {
			return true
		}
		
	}
	return false
}

// remove an existing volume
func VolumeRemove(name string) error {
	return client.VolumeRemove(context.Background(), name, true)
}