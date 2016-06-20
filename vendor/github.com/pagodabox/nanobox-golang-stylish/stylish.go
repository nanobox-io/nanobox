// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

package stylish

import (
	"fmt"
	"strings"

	"github.com/mitchellh/go-wordwrap"
)


// Nest is a generic nesting function that 
// will generate the appropariate prefix based
// on the nest level
func Nest(level int, msg string) (rtn string) {
	for index, line := range strings.Split(msg, "\n") {
		// skip the last new line at the end of the message
		// because we add the new line in on each Sprintf
		if index == len(strings.Split(msg, "\n")) - 1 && line == "" {
			continue
		}
		rtn += fmt.Sprintf("%s%s\n", GenerateNestedPrefix(level), shorten(level, line))
	}
	return
}

// ProcessStart styles and prints a 'child process' as outlined at:
// http://nanodocs.gopagoda.io/engines/style-guide#child-process
//
// Usage:
// ProcessStart "i am a process"
//
// Output:
// + I am a process ------------------------------------------------------------ >
func ProcessStart(msg string, v ...interface{}) string {

	maxLen := 80
	process := fmt.Sprintf(msg, v...)
	subLen := len(process) + len("+   >")

	// print process, inserting a '-' (colon) 'n' times, where 'n' is the number
	// remaining after subtracting subLen (number of 'reserved' characters) from
	// maxLen (maximum number of allowed characters)
	return fmt.Sprintf("+ %s %s >\n", process, strings.Repeat("-", (maxLen-subLen)))
}

// ProcessEnd styles and prints a 'child process' as outlined at:
// http://nanodocs.gopagoda.io/engines/style-guide#child-process
//
// Usage:
// ProcessEnd
//
// Output:
// <new line>
func ProcessEnd() string {
	return fmt.Sprintf("\n")
}

// Marker is the root for Bullet/SubBullet; used alone, it allows for a custom
// mark to be specified
//
// Usage:
// Maker "*",  "i am a marker"
//
// Output:
// * i am a marker
func Marker(mark, msg string, v ...interface{}) string {
	return fmt.Sprintf("%s %s\n", mark, fmt.Sprintf(msg, v...))
}

// Bullet styles and prints a message as outlined at:
// http://nanodocs.gopagoda.io/engines/style-guide#bullet-points
//
// Usage:
// Bullet "i am a bullet"
//
// Output:
// + i am a bullet
func Bullet(msg string, v ...interface{}) string {
	return Marker("+", fmt.Sprintf(msg, v...))
}

// SubBullet styles and prints a message as outlined at:
// http://nanodocs.gopagoda.io/engines/style-guide#bullet-points
//
// Usage:
// SubBullet "i am a sub bullet"
//
// Output:
//   i am a sub bullet
func SubBullet(msg string, v ...interface{}) string {
	return Marker("  +", fmt.Sprintf(msg, v...))
}

// Warning styles and prints a message as outlined at:
// http://nanodocs.gopagoda.io/engines/style-guide#warning
//
// Usage:
// Warning "You just bought Hot Pockets!"
//
// Output:
// -----------------------------  WARNING  -----------------------------
// You just bought Hot Pockets!
func Warning(body string, v ...interface{}) string {
	return fmt.Sprintf(`
----------------------------------  WARNING  ----------------------------------
%s
`, wordwrap.WrapString(fmt.Sprintf(body, v...), 70))
}

// ErrorHead styles and prints an error heading as outlined at:
// http://nanodocs.gopagoda.io/engines/style-guide#fatal_errors
//
// Usage:
// ErrorHead "nuclear launch detected"
//
// Output:
// ! NUCLEAR LAUNCH DETECTED !
func ErrorHead(heading string, v ...interface{}) string {
	return fmt.Sprintf("\n! %s !\n", strings.ToUpper(fmt.Sprintf(heading, v...)))
}

// ErrorBody styles and prints an error body as outlined at:
// http://nanodocs.gopagoda.io/engines/style-guide#fatal_errors
//
// Usage:
// ErrorBody "All your base are belong to us"
//
// Output:
// All your base are belong to us
func ErrorBody(body string, v ...interface{}) string {
	return fmt.Sprintf("%s\n", wordwrap.WrapString(fmt.Sprintf(body, v...), 70))
}

// Error styles and prints a message as outlined at:
// http://nanodocs.gopagoda.io/engines/style-guide#fatal_errors
//
// Usage:
// Error "nuclear launch detected", "All your base are belong to us"
//
// Output:
// ! NUCLEAR LAUNCH DETECTED !
//
// All your base are belong to us
func Error(heading, body string) string {
	return fmt.Sprintf("%s%s", ErrorHead(heading), ErrorBody(body))
}

// GenerateNestedPrefix will generate a prefix string of spaces to match the
// specified depth level
func GenerateNestedPrefix(level int) string {
	prefix := ""

	for i := 0; i < level; i++ {
		prefix += "  "
	}

	return prefix
}

func shorten(level int, msg string) string {
	switch {
	case isProgress(msg):
		return shortenProgress(level, msg)
	case isWarning(msg):
		return shortenWarning(level, msg)
	}
	return msg
}

func isProgress(line string) bool {
	return len(line) == 80 && strings.HasSuffix(line, "-- >")
}

func shortenProgress(level int, msg string) string {
	suffix := fmt.Sprintf("%s >", strings.Repeat("-", (level * 2)))
	return strings.Replace(msg, suffix, " >", 1)
}

func isWarning(line string) bool {
	return line == "----------------------------------  WARNING  ----------------------------------"
}

func shortenWarning(level int, msg string) string {
	wrapper := strings.Repeat("-", level)
	return strings.Replace(msg, fmt.Sprintf("%s  WARNING  %s", wrapper, wrapper), "  WARNING  ", 1)
}
