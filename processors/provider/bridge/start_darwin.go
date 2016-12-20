package bridge

import (
	"fmt"

	"os/exec"
	"io/ioutil"
)

func serviceConfigFile() string {
	return "/Library/LaunchDaemons/io.nanobox.openvpn.plist"
}

func serviceConfig() string {
	return fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
        <key>label</key>
        <string>io.nanobox.openvpn</string>

        <key>ProgramArguments</key>
        <array>
                <string>%s</string>
                <string>--config</string>
                <string>%s</string>
        </array>
</dict>
</plist>
`, bridgeClient, configFile())
}

func createService() error {

	// setup config file
	if err := ioutil.WriteFile(serviceConfigFile(), []byte(serviceConfig()), 0644); err != nil {
		return err
	}

	out, err := exec.Command("launchctl","load", "-w", "/Library/LaunchDaemons/io.nanobox.openvpn.plist").CombinedOutput() 
	if err != nil {
		fmt.Errorf("out: %s, err: %s", out, err)
	}
	return nil
}

func startService() error {
	out, err := exec.Command("launchctl", "start", "io.nanobox.openvpn").CombinedOutput() 
	if err != nil {
		fmt.Errorf("out: %s, err: %s", out, err)
	}
	return nil
}
