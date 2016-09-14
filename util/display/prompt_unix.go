package display

import (
	"fmt"
	"os"
	"golang.org/x/crypto/ssh/terminal"
)

// ReadPassword reads a password from the terminal and masks the input
func ReadPassword() (string, error) {
	fmt.Print("Password: ")

	pass, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println("")

	return string(pass), err
}