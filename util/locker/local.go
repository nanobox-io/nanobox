package locker

import (
	"fmt"
	"net"
	"time"

	"github.com/jcelliott/lumber"
	"github.com/nanobox-io/nanobox/util/config"
)

var (
	lln    net.Listener // local locking network
	lCount int
)

// LocalLock locks on port
func LocalLock() error {

	//
	for {
		if success, _ := LocalTryLock(); success {
			break
		}
		lumber.Trace("local lock waiting...")
		<-time.After(time.Second)
	}

	mutex.Lock()
	lCount++
	lumber.Trace("local lock aquired (%d)", lCount)
	mutex.Unlock()

	return nil
}

// LocalTryLock ...
func LocalTryLock() (bool, error) {

	var err error

	//
	if lln != nil {
		return true, nil
	}

	//
	port := config.Viper().GetInt("lock-port")
	if port == 0 {
		port = 12345
	}
	port = port + localPort()

	//
	if lln, err = net.Listen("tcp", fmt.Sprintf(":%d", port)); err == nil {
		return true, nil
	}

	return false, nil
}

// LocalUnlock ...
func LocalUnlock() (err error) {

	mutex.Lock()
	lCount--
	lumber.Trace("local lock released (%d)", lCount)
	mutex.Unlock()

	// if im not the last guy to release my lock quit immidiately instead of closing
	// the connection
	if lCount > 0 || lln == nil {
		return nil
	}

	err = lln.Close()
	lln = nil

	return
}

// localPort ...
func localPort() (num int) {

	b := []byte(config.AppID())

	//
	for i := 0; i < len(b); i++ {
		num = num + int(b[i])
	}

	return num
}
