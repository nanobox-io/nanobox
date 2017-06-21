// +build !windows

package display

import (
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"os"
)

// ReadPassword reads a password from the terminal and masks the input
func ReadPassword() (string, error) {
	fmt.Print("Nanobox Password: ")

	pass, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println("")

	return string(pass), err
}
