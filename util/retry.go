package util

import (
	"time"
)

type Retryable func() error

func Retry(retryFunc Retryable, attempts int, delay time.Duration) (err error) {

	for i := 0; i < attempts; i++ {
		err = retryFunc()
		if err == nil {
			return
		}
		// delay
		<-time.After(delay)
	}
	return
}