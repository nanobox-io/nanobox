package service

import (
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr"
)

func Running(name string) bool {

	m, err := mgr.Connect()
	if err != nil {
		return false
	}
	defer m.Disconnect()

	// check to see if we need to create at all
	s, err := m.OpenService(name)
	if err != nil {
		// jobs done
		return false
	}
	defer s.Close()

	status, err := s.Query()
	if err != nil {
		return false
	}

	return status.State == svc.Running
}
