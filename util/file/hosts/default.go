// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

//
package hosts

type (
	host struct{}
	Host interface{}
)

var (
	Default Host = host{}
)

func (host) HasDomain() bool {
	return HasDomain()
}

func (host) AddDomain() {
	AddDomain()
}

func (host) RemoveDomain() {
	RemoveDomain()
}
