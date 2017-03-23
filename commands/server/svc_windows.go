package server

import (
	"fmt"
	"os"
	"time"

	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/debug"
	"golang.org/x/sys/windows/svc/eventlog"
	
)

var elog debug.Log

type nanoboxServer struct {}

func (ns *nanoboxServer) Execute(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (ssec bool, errno uint32) {
	elog.Info(1, fmt.Sprintf("execute called"))
	const cmdsAccepted = svc.AcceptStop | svc.AcceptShutdown
	changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}
	elog.Info(1, fmt.Sprintf("running"))

loop:
	for {
		select {
		case c := <-r:
			switch c.Cmd {
			case svc.Interrogate:
				changes <- c.CurrentStatus
				// Testing deadlock from https://code.google.com/p/winsvc/issues/detail?id=4
				time.Sleep(100 * time.Millisecond)
				changes <- c.CurrentStatus
			case svc.Stop, svc.Shutdown:
				break loop
			default:
				elog.Error(1, fmt.Sprintf("unexpected control request #%#v", c))
			}
		}
	}
	changes <- svc.Status{State: svc.StopPending}
	return
}

func svcStart() {
	var err error
	elog, err = eventlog.Open(name)
	if err != nil {
		return
	}

	defer elog.Close()

	elog.Info(1, fmt.Sprintf("starting %s service", name))
	err = svc.Run(name, &nanoboxServer{})
	if err != nil {
		elog.Error(1, fmt.Sprintf("%s service failed: %v", name, err))
		return
	}
	elog.Info(1, fmt.Sprintf("%s service stopped", name))
	
	// on windows when i get here the service has been stopped so we exit
	os.Exit(0)
}