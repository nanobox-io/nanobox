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

// Header styles and prints a header as outlined at:
// http://nanodocs.gopagoda.io/engines/style-guide#header
//
// Usage:
// Header "i am a header"
//
// Output:
// :::::::::::::::::::::::::: I AM A HEADER :::::::::::::::::::::::::
func Header(msg string, v ...interface{}) string {

	maxLen := 70
	subLen := len(fmt.Sprintf(msg, v...))

	leftLen := (maxLen-subLen)/2 + (maxLen-subLen)%2
	rightLen := (maxLen - subLen) / 2

	heading := strings.ToUpper(fmt.Sprintf(msg, v...))

	// print heading, inserting a ':' (colon) 'n' times, where 'n' is the number
	// remaining after subtracting subLen (number of 'reserved' characters) from
	// maxLen (maximum number of allowed characters)
	return fmt.Sprintf(`
%s
`, fmt.Sprintf("%s %s %s", strings.Repeat(":", leftLen), heading, strings.Repeat(":", rightLen)))
}

// ProcessStart styles and prints a 'child process' as outlined at:
// http://nanodocs.gopagoda.io/engines/style-guide#child-process
//
// Usage:
// ProcessStart "i am a process"
//
// Output:
// I AM A PROCESS -------------------------------------------------->
func ProcessStart(msg string, v ...interface{}) string {

	maxLen := 70
	subLen := len(fmt.Sprintf("%s%s->", fmt.Sprintf(msg, v...)))

	process := strings.ToUpper(fmt.Sprintf(msg, v...))

	// print process, inserting a '-' (colon) 'n' times, where 'n' is the number
	// remaining after subtracting subLen (number of 'reserved' characters) from
	// maxLen (maximum number of allowed characters)
	return fmt.Sprintf(`
%s
`, fmt.Sprintf("%s %s->", process, strings.Repeat("-", (maxLen-subLen))))
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

// SubTask styles and prints a 'sub task' as outlined at:
// http://nanodocs.gopagoda.io/engines/style-guide#sub-tasks
//
// Usage:
// SubTask "i am a sub task"
//
// Output:
// I AM A SUB TASK ----------------------->
func SubTaskStart(msg string, v ...interface{}) string {

	maxLen := 40
	subLen := len(fmt.Sprintf("%s->", fmt.Sprintf(msg, v...)))

	task := strings.ToUpper(fmt.Sprintf(msg, v...))

	// print msg, inserting a ':' (colon) 'n' times, where 'n' is the number
	// remaining after subtracting subLen (number of 'reserved' characters) from
	// maxLen (maximum number of allowed characters)
	return fmt.Sprintf(`
%s
`, fmt.Sprintf("%s %s->", task, strings.Repeat("-", (maxLen-subLen))))
}

// SubTaskSuccess styles and prints a footer to a successful subtask
//
// Usage:
// SubTaskSuccess
//
// Output:
//    [√] SUCCESS
func Success() string {
	return fmt.Sprintf("   [√] SUCCESS\n")
}

// SubTaskFail styles and prints a footer to a failed subtask
//
// Usage:
// SubTaskFail
//
// Output:
//    [!] FAILED
func Fail() string {
	return fmt.Sprintf("   [!] FAILED\n")
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
// +> i am a bullet
func Bullet(msg string, v ...interface{}) string {
	return Marker("+>", fmt.Sprintf(msg, v...))
}

// ErrBullet is a Bullet to be used for quick inline failure messsages
//
// Usage:
// ErrBullet "i am an errBullet"
//
// Output:
// -> i am an errBullet
func ErrBullet(msg string, v ...interface{}) string {
	return Marker("->", fmt.Sprintf(msg, v...))
}

// SubBullet styles and prints a message as outlined at:
// http://nanodocs.gopagoda.io/engines/style-guide#bullet-points
//
// Usage:
// SubBullet "i am a sub bullet"
//
// Output:
//    i am a sub bullet
func SubBullet(msg string, v ...interface{}) string {
	return Marker("  ", fmt.Sprintf(msg, v...))
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
-----------------------------  WARNING  -----------------------------
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
