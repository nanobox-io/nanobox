package service

import (
	"fmt"
	"net"
	"time"
)

func serviceConfigFile(name string) string {
	return fmt.Sprintf("/Library/LaunchDaemons/io.%s.plist", name)
}

func startCmd(name string) []string {
	return []string{"launchctl", "start", fmt.Sprintf("io.%s", name)}
}

func Running(name string) bool {
	<-time.After(500*time.Millisecond)
	conn, err := net.DialTimeout("tcp", "127.0.0.1:23456", 100*time.Millisecond)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

func stopCmd(name string) []string {
	return []string{"launchctl", "stop", fmt.Sprintf("io.%s", name)}
}

func removeCmd(name string) []string {
	return []string{"launchctl", "remove", fmt.Sprintf("io.%s", name)}
}
