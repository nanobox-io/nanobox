// +build !windows

package ui

import (
	"code.google.com/p/gopass"
	"fmt"
)

// PPrompt prompts for a password but keeps the typed response hidden
func PPrompt(p string) string {
	password, err := gopass.GetPass(p)
	if err != nil {
		fmt.Println("Unable to read input. See ~/.pagodabox/log.txt for details")
		Error("ui.PPrompt", err)
	}

	return password
}
