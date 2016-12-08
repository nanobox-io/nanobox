package util_test

import (
	"fmt"
	"time"
	"testing"

	"github.com/nanobox-io/nanobox/util"
)

func TestRetry(t *testing.T) {
	failureCount := 0
	failingFunc := func() error {
		failureCount += 1
		if failureCount > 5 {
			return nil
		}
		return fmt.Errorf("error")
	}

	err := util.Retry(failingFunc, 3, time.Nanosecond)
	if err == nil {
		t.Errorf("func failed but didnt error")
	}

	err = util.Retry(failingFunc, 3, time.Nanosecond)
	if err != nil {
		t.Errorf("func succeeded but i recieved an error")
	}


}