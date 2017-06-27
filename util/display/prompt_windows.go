// +build windows

package display

import (
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"os"
)

// ReadPassword reads a password from the terminal and masks the input
func ReadPassword(label string) (string, error) {

	// Fetch the current state of the terminal so it can be restored later
	oldState, err := terminal.GetState(int(os.Stdin.Fd()))
	if err != nil {
		return "", err
	}
	// Turn off echo and make stdin blank
	terminal.MakeRaw(int(os.Stdin.Fd()))
	// Restore echo after the function exits
	defer terminal.Restore(int(os.Stdin.Fd()), oldState)

	fmt.Printf("%s Password: ", label)

	// Read the password from stdin
	t := terminal.NewTerminal(os.Stdin, "")
	pass, err := t.ReadPassword("")

	// Add a newline so the next output isn't next to the Password:
	fmt.Println("")

	if err != nil {
		return "", err
	}

	return pass, nil
}
