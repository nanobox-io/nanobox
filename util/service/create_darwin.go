package service

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"
)

func Create(name string, command []string) error {

	// setup config file
	if err := ioutil.WriteFile(serviceConfigFile(name), []byte(serviceConfig(name, command)), 0644); err != nil {
		return err
	}

	out, err := exec.Command("launchctl", "load", serviceConfigFile(name)).CombinedOutput()
	if err != nil {
		fmt.Errorf("out: %s, err: %s", out, err)
	}

	return nil
}

func serviceConfig(name string, command []string) string {
	return fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
        <key>Label</key>
        <string>io.%s</string>

        <key>ProgramArguments</key>
        <array>
                <string>%s</string>
        </array>
</dict>
</plist>
`, name, strings.Join(command, "</string>\n                <string>"))
}
