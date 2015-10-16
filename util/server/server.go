// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

//
package server

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/nanobox-io/nanobox/config"
)

var (
	err error //
)

// Service
type Service struct {
	CreatedAt time.Time
	EnvVars   map[string]string
	IP        string
	Name      string
	Password  string
	Ports     []int
	Username  string
	UID       string
}

// Get issues a "GET" request to nanobox server at the specified path
func Get(path string, v interface{}) (*http.Response, error) {

	req, err := http.NewRequest("GET", config.ServerURL+path, nil)
	if err != nil {
		return nil, err
	}

	//
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	// read the body
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	// unmarshal response into the provided container. If no container was given
	// it's mostly likely a raw request where the body isn't needed, and therfore
	// this step can be skipped.
	if v != nil {
		if err := json.Unmarshal(b, &v); err != nil {
			return nil, err
		}
	}

	return res, nil
}

// Post issues a "POST" to nanobox server
func Post(path, contentType string, body io.Reader) (*http.Response, error) {

	res, err := http.Post(config.ServerURL+path, contentType, body)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// Put
func Put(path string, body io.Reader) (*http.Response, error) {

	// check to see if the VM is able to be suspended
	req, err := http.NewRequest("PUT", config.ServerURL+path, body)
	if err != nil {
		return nil, err
	}

	//
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
