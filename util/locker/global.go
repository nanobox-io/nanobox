package locker

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/jcelliott/lumber"

	"github.com/nanobox-io/nanobox/models"
)

var (
	gln    net.Listener
	gCount int
	mutex  = sync.Mutex{}
)

// GlobalLock locks on port
func GlobalLock() error {
	lumber.Trace("global locking")

	//
	for {
		if success, _ := GlobalTryLock(); success {
			break
		}
		lumber.Trace("global lock waiting...")
		<-time.After(time.Second)
	}

	mutex.Lock()
	gCount++
	lumber.Trace("global lock aqquired (%d)", gCount)
	mutex.Unlock()

	return nil
}

// GlobalTryLock ...
func GlobalTryLock() (bool, error) {

	var err error

	//
	if gln != nil {
		return true, nil
	}

	//
	config, _ := models.LoadConfig()
	port := config.LockPort
	if port == 0 {
		port = 12345
	}

	//
	if gln, err = net.Listen("tcp", fmt.Sprintf(":%d", port)); err == nil {
		return true, nil
	}

	return false, nil
}

// GlobalUnlock removes the lock if im the last global unlock to be called; this
// needs to be called EXACTLY he same number of tiems as lock
func GlobalUnlock() error {

	mutex.Lock()
	gCount--
	lumber.Trace("global lock released (%d)", gCount)
	mutex.Unlock()

	//
	if gCount > 0 || gln == nil {
		return nil
	}

	err := gln.Close()
	gln = nil

	return err
}
