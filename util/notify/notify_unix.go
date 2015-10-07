// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//
// +build linux darwin

//
package notify

import "syscall"

// setFileLimit sets rlimit to the maximum rlimit
func setFileLimit() {
	rlm := &syscall.Rlimit{}
	syscall.Getrlimit(syscall.RLIMIT_NOFILE, rlm)
	rlm.Cur = rlm.Max
	syscall.Setrlimit(syscall.RLIMIT_NOFILE, rlm)
}
