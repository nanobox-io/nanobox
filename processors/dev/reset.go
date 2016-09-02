package dev

import (
	"fmt"
	
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util/display"
)

// Reset all the counters for all dev applications
func Reset() error {
	display.OpenContext("Resetting dev state")
	defer display.CloseContext()
	
	display.StartTask("Clear usage counters")
	if err := models.DeleteAllCounters(); err != nil {
		display.ErrorTask()
		return fmt.Errorf("failed to remove usage counters: %s", err.Error())
	}
	display.StopTask()
	
	return nil
}
