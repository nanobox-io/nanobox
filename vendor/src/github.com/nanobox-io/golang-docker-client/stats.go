package docker

import (
	"encoding/json"

	dockType "github.com/docker/engine-api/types"
	"golang.org/x/net/context"
)

// ContainerStats(ctx context.Context, containerID string, stream bool) (io.ReadCloser, error)
func ContainerStats(id string) (dockType.Stats, error) {
	rc, err := client.ContainerStats(context.Background(), id, false)
	var stats dockType.Stats
	decoder := json.NewDecoder(rc)
	for decoder.More() {
		decoder.Decode(&stats)
		// fmt.Printf("STATS!!!%+v\n", stats)
	}
	defer rc.Close()
	return stats, err
}
