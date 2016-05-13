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
	subLen := len(fmt.Sprintf("+ %s%s >", fmt.Sprintf(msg, v...)))

	process := fmt.Sprintf(msg, v...)

	// print process, inserting a '-' (colon) 'n' times, where 'n' is the number
	// remaining after subtracting subLen (number of 'reserved' characters) from
	// maxLen (maximum number of allowed characters)
	return fmt.Sprintf("%s\n", fmt.Sprintf("+ %s %s >", process, strings.Repeat("-", (maxLen-subLen))))
}

// NestedProcessStart styles and prints a 'child process' as outlined at:
// http://nanodocs.gopagoda.io/engines/style-guide#child-process
// with a nested prefix according to the level specified
func NestedProcessStart(msg string, level int) string {
	return fmt.Sprintf("%s%s", GenerateNestedPrefix(level), ProcessStart(msg))
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

// NestedBullet styles and prints a message as outlined at:
// http://nanodocs.gopagoda.io/engines/style-guide#bullet-points
// with a nested prefix according to the level specified
func NestedBullet(msg string, level int) string {
	return fmt.Sprintf("%s%s", GenerateNestedPrefix(level), Bullet(msg))
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
	return Marker(" ", fmt.Sprintf(msg, v...))
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
	return fmt.Sprintf("%s\n", wordwrap.WrapString(fmt.Sprintf(body, v...), 80))
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
