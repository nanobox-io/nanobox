package dev

import (
	"github.com/nanobox-io/nanobox/models"
)

// Reset all the counters for all dev applications
type Reset struct {
}

//
func (reset Reset) Run() error {
	return models.DeleteAllCounters()	
}
