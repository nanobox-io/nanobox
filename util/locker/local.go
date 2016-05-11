package locker

import (
	"fmt"
	"net"
	"time"

	"github.com/jcelliott/lumber"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/nanofile"
)

// local locking network
var lln net.Listener
var lCount int = 0

// Lock locks on port
func LocalLock() error {
	for {
		if success, _ := LocalTryLock(); success {
			break
		}
		lumber.Debug("local lock waiting...")
		<-time.After(time.Second)
	}
	mutex.Lock()
	lCount++
	lumber.Debug("local lock aquired (%d)", lCount)
	mutex.Unlock()
	return nil
}

func LocalTryLock() (bool, error) {
	if lln != nil {
		return true, nil
	}
	var err error
	port := nanofile.Viper().GetInt("lock-port")
	if port == 0 {
		port = 12345
	}
	port = port + localPort()

	if lln, err = net.Listen("tcp", fmt.Sprintf(":%d", port)); err == nil {
		return true, nil
	}
	return false, nil
}

func LocalUnlock() (err error) {
	mutex.Lock()
	lCount--
	lumber.Debug("local lock released (%d)", lCount)	
	mutex.Unlock()
	// if im not the last guy to release my lock
	// quit immidiately instead of closing the connection
	if lCount > 0 || lln == nil {
		return nil
	}
	err = lln.Close()
	lln = nil
	return
}

func localPort() int {
	b := []byte(util.AppName())
	num := 0
	for i := 0; i < len(b); i++ {
		num = num + int(b[i])
	}
	return num
}
