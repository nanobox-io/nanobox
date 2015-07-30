// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

// package hatchet provides a very simple Logger interface that is intentionally
// generic allowing it to be used with numerous loggers that are already in existance,
// or the easy creation of a custom logger. It also provides a DevNullLogger which
// can be used when a project requires a logger but none is provided.
package hatchet

//
type (

	// Logger is a simple interface that's designed to be intentionally generic to
	// allow many different types of Logger's to satisfy its interface
	Logger interface {
		Fatal(string, ...interface{})
		Error(string, ...interface{})
		Warn(string, ...interface{})
		Info(string, ...interface{})
		Debug(string, ...interface{})
		Trace(string, ...interface{})
	}

	// DevNullLogger is a Logger that purpose is to provide no output
	DevNullLogger struct{}
)

// The following methods are provided on DevNullLogger so that it satisfies the
// Logger interface, but allows it to log nothing.
func (d DevNullLogger) Fatal(s string, v ...interface{}) {}
func (d DevNullLogger) Error(s string, v ...interface{}) {}
func (d DevNullLogger) Warn(s string, v ...interface{})  {}
func (d DevNullLogger) Info(s string, v ...interface{})  {}
func (d DevNullLogger) Debug(s string, v ...interface{}) {}
func (d DevNullLogger) Trace(s string, v ...interface{}) {}
