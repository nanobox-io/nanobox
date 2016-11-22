package docker

import (
	"encoding/json"

	dockType "github.com/docker/engine-api/types"
	"golang.org/x/net/context"
)

// ContainerStats(ctx context.Context, containerID string, stream bool) (io.ReadCloser, error)
func ContainerStats(id string) (dockType.Stats, error) {
	rc, err := client.ContainerStats(context.Background(), id, false)
	if err != nil {
		return dockType.Stats{}, err
	}
	defer rc.Close()

	var stats dockType.Stats
	decoder := json.NewDecoder(rc)
	for decoder.More() {
		decoder.Decode(&stats)
	}

	return stats, err
}
