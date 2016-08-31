package dev

import (
	"github.com/nanobox-io/nanobox/models"
)

// Reset all the counters for all dev applications
func Reset() error {
	return models.DeleteAllCounters()
}
