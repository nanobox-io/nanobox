// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

package client

import (
	"net/url"
	"time"
)

//
type (

	// User represents a nanobox user
	User struct {
		AuthenticationToken string     `json:"authentication_token"` //
		CreatedAt           time.Time  `json:"created_at"`           //
		Email               string     `json:"email"`                //
		ID                  string     `json:"id"`                   //
		UpdatedAt           time.Time  `json:"updated_at"`           //
		Username            string     `json:"username"`             //
	}

	// UserUpdateOptions represents all available options when updating a user
	UserUpdateOptions struct{}
)

// GetAuthToken takes a userSlug and password to return a user's authentication
// token
func GetAuthToken(userSlug, password string) (*User, error) {

	//
	v := url.Values{}
	v.Set("id", userSlug)
	v.Add("password", password)

	// this path is used (vs restful) to avoid sending emails as part of the path
	path := APIURL + "/" + APIVersion + "/user_auth_token?" + v.Encode()

	var user User
	return &user, DoRawRequest(&user, "GET", path, nil, nil)
}
