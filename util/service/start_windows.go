package service

import (
	"fmt"
	"strings"

	"golang.org/x/sys/windows/svc/mgr"
)

func Start(name string) error {
	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()

	s, err := m.OpenService(name)
	if err != nil {
		return fmt.Errorf("could not access service: %v", err)
	}
	defer s.Close()

	err = s.Start()
	if err != nil {
		if strings.Contains(err.Error(), "already running") {
			return nil
		}
		return fmt.Errorf("could not start service: %v", err)
	}

	return nil
}
