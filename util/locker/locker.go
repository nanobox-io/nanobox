package locker

import (
	"fmt"
	"net"
	"time"

	"github.com/nanobox-io/nanobox/util/nanofile"
)

var ln net.Listener

// Lock locks on port
func Lock() error {
	for {
		if ok, err := TryLock(); ok {
			return err
		}
		<-time.After(time.Second)
	}
	return nil
}

func TryLock() (bool, error) {
	if ln != nil {
		return true, nil
	}
	var err error
	port := nanofile.Viper().GetInt("lock-port")
	if port == 0 {
		port = 12345
	}
	if ln, err = net.Listen("tcp", fmt.Sprintf(":%d", port)); err == nil {
		return true, nil
	}
	return false, nil
}

func Unlock() error {
	if ln == nil {
		return nil
	}
	err := ln.Close()
	ln = nil
	return err
}
