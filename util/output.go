// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

package util

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/howeyc/gopass"
	"github.com/mitchellh/colorstring"

	"github.com/pagodabox/nanobox-cli/config"
)

// Printv (print verbose) only prints a message if the 'verbose' flag is passed
func Printv(msg string, verbose bool) {
	if verbose {
		fmt.Printf(msg)
	}
}

// Printc wraps a print message in 'colorstring' and passes it to fmt.Println
func Printc(msg string, v ...interface{}) {
	fmt.Println(colorstring.Color(fmt.Sprintf(msg, v...)))
}

// Debug
func Debug(msg string, debug bool) {
	if debug {
		fmt.Printf(msg)
	}
}

// Log
func Log(msg string, debug bool) {
	config.Log.Info(msg)
}

// Fatal
func Fatal(msg string, err error) {
	fmt.Println("Fatal error! See ~/.nanobox/nanobox.log for details. Exiting...")
	config.Log.Fatal(fmt.Sprintf("%v - %v", msg, err))
	config.Log.Close()
	os.Exit(1)
}

// Prompt will prompt for input from the shell and return a trimmed response
func Prompt(p string, v ...interface{}) string {
	reader := bufio.NewReader(os.Stdin)

	//
	fmt.Print(colorstring.Color(fmt.Sprintf(p, v...)))

	input, err := reader.ReadString('\n')
	if err != nil {
		Fatal("[utils/ui] reader.ReadString() failed", err)
	}

	return strings.TrimSpace(input)
}

// PPrompt prompts for a password but keeps the typed response hidden
func PPrompt(p string) string {
	fmt.Printf(p)
	return string(gopass.GetPasswd())
}
