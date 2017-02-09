package util_test

import (
	"fmt"
	"testing"
	"time"

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

func TestError(t *testing.T) {
	err := util.ErrorfQuiet("hi %s", "world")
	if err.Error() != "hi world" {
		t.Errorf("did not format correctly")
	}

	err = util.ErrorAppend(err, "james")
	if err.Error() != "james: hi world" {
		t.Errorf("append failed")
	}
}
