// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

// +build !windows

package ui

import "code.google.com/p/gopass"

// PPrompt prompts for a password but keeps the typed response hidden
func PPrompt(p string) string {
	password, err := gopass.GetPass(p)
	if err != nil {
		LogFatal("[ui.ui_unix] PPrompt() failed", err)
	}

	return password
}
