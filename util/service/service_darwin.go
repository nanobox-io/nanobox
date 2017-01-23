package service

import "fmt"

func serviceConfigFile(name string) string {
	return fmt.Sprintf("/Library/LaunchDaemons/io.%s.plist", name)
}

func startCmd(name string) []string {
	return []string{"launchctl", "start", fmt.Sprintf("io.%s", name)}
}

func stopCmd(name string) []string {
	return []string{"launchctl", "stop", fmt.Sprintf("io.%s", name)}	
}

func removeCmd(name string) []string {
	return []string{"launchctl", "remove", fmt.Sprintf("io.%s", name)}
}