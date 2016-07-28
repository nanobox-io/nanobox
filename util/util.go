// Package util ...
package util

import (
	"fmt"
	"math/rand"
	"os"
	"time"
	
	"golang.org/x/crypto/ssh/terminal"
)

const (

	// VERSION is the global version for nanobox; mainly used in the update process
	// but placed here to allow access when needed (commands, processor, etc.)
	VERSION     = "1.0.0"
	letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

// RandomString ...
func RandomString(size int) string {

	// create a new randomizer with a unique seed
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	//
	b := make([]byte, size)
	for i := range b {
		b[i] = letterBytes[r.Intn(len(letterBytes))]
	}

	return string(b)
}

// ReadPassword reads a password from the terminal and masks the input
func ReadPassword() (string, error) {
	
	// Fetch the current state of the terminal so it can be restored later
	oldState, err := terminal.GetState(int(os.Stdin.Fd()))
	if err != nil {
		return "", err
	}
	// Turn off echo and make stdin blank
	terminal.MakeRaw(int(os.Stdin.Fd()))
	// Restore echo after the function exits
	defer terminal.Restore(int(os.Stdin.Fd()), oldState)
	
	fmt.Printf("Password: ")

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
