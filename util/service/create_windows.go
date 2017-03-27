package service

import (
	// "fmt"
	// // "io/ioutil"
	// "os/exec"
	// "strings"

	"golang.org/x/sys/windows/svc/eventlog"
	"golang.org/x/sys/windows/svc/mgr"
)

func Create(name string, command []string) error {
	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()

	// check to see if we need to create at all
	s, err := m.OpenService(name)
	if err == nil {
		s.Close()
		// jobs done
		return nil
	}

	// create the service
	args := []string{}
	if len(command) > 1 {	
		args = command[1:]
	}
	s, err = m.CreateService(name, command[0], mgr.Config{DisplayName: name}, args...)
	if err != nil {
		return err
	}
	defer s.Close()

	err = eventlog.InstallAsEventCreate(name, eventlog.Error|eventlog.Warning|eventlog.Info)
	if err != nil {
		// s.Delete()
		// eventlog.Remove(name)
		// return fmt.Errorf("SetupEventLogSource() failed: %s", err)
	}
	return nil
}
