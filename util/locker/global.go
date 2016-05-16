package locker

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/jcelliott/lumber"

	"github.com/nanobox-io/nanobox/util/nanofile"
)

var gln net.Listener
var gCount int = 0
var mutex = sync.Mutex{}

// Lock locks on port
func GlobalLock() error {
	lumber.Debug("global locking")
	for {
		if success, _ := GlobalTryLock(); success {
			break
		}
		lumber.Debug("global lock waiting...")
		<-time.After(time.Second)
	}

	mutex.Lock()
	gCount++
	lumber.Debug("global lock aqquired (%d)", gCount)
	mutex.Unlock()
	return nil
}

func GlobalTryLock() (bool, error) {
	if gln != nil {
		return true, nil
	}
	var err error
	port := nanofile.Viper().GetInt("lock-port")
	if port == 0 {
		port = 12345
	}
	if gln, err = net.Listen("tcp", fmt.Sprintf(":%d", port)); err == nil {
		return true, nil
	}
	return false, nil
}

// remove the lock if im the last global unlock to be called
// this needs to be called exactlyt he same number of tiems as lock
func GlobalUnlock() error {
	mutex.Lock()
	gCount--
	lumber.Debug("global lock released (%d)", gCount)
	mutex.Unlock()
	if gCount > 0 || gln == nil {
		return nil
	}
	err := gln.Close()
	gln = nil
	return err
}
