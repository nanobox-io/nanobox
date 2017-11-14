// Package cache contains logic for managing the cache volume.
package cache

import (
	"fmt"

	"github.com/nanobox-io/golang-docker-client"

	"github.com/nanobox-io/nanobox/models"
)

// ClearAll will clear all cache volumes for the app environments passed in.
func ClearAll(envs []*models.Env) error {
	for i := range envs {
		err := docker.VolumeRemove(fmt.Sprintf("nanobox_%s_cache", envs[i].ID))
		if err != nil {
			return fmt.Errorf("Failed to remove volume - %s", err.Error())
		}
	}

	return nil
}

// Clear will clear the cache volume for the app id passed in.
func Clear(id string) error {
	err := docker.VolumeRemove(fmt.Sprintf("nanobox_%s_cache", id))
	if err != nil {
		return fmt.Errorf("Failed to remove volume - %s", err.Error())
	}

	return nil
}
