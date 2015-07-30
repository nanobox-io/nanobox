// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

// package client consists of a core api client struct with methods broken into
// related calls, for interacting and communicating with the nanobox API.
package client

//
import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
)

//
const (
	DefaultAPIURL      = "https://api.nanobox.io"
	DefaultAPIVersion  = "v1"
	DefaultContentType = "application/json"
	Version            = "1.0.0"
)

//
var (
	APIURL     string       // The URL to which the API client will make requests.
	APIVersion string       // The version of the API to make requests to.
	AuthToken  string       // The authentication_token of the user to make requests with.
	Debug      bool         // If debug mode is enabled.
	HTTPClient *http.Client // The HTTP.Client to use when making requests.
	UserSlug   string       // The UserSlug to use in conjunction with the AuthToken when making API requests. (username, email, or ID)
)

//
type (

	// APIError represents a pagoda-client error
	APIError struct {
		error         // The entire error (ex. {"error":"404 Not Found"})
		Body   string `json:"error"` // The error body (ex. "Not Found")
		Code   int    // The 'int' status code (ex. 404)
		Status string `json:"status"` // The 'string' status code (ex. "404")
	}

	// Email represents an email that can be attached to objects like cron jobs or
	// invoices
	Email struct {
		Email string
	}
)

//
func init() {
	APIURL = DefaultAPIURL
	APIVersion = DefaultAPIVersion
	Debug = false
	HTTPClient = http.DefaultClient
}

// post handles standard POST operations to the nanobox API
func post(v interface{}, path string, body interface{}) error {
	return doAPIRequest(v, "POST", path, body)
}

// get handles standard GET operations to the nanobox API
func get(v interface{}, path string) error {
	return doAPIRequest(v, "GET", path, nil)
}

// patch handles standard PATH operations to the nanobox API
func patch(v interface{}, path string, body interface{}) error {
	return doAPIRequest(v, "PATCH", path, body)
}

// put handles standard PUT operations to the nanobox API
func put(v interface{}, path string, body interface{}) error {
	return doAPIRequest(v, "PUT", path, body)
}

// delete handles standard DELETE operations to the nanobox API
func delete(path string) error {
	return doAPIRequest(nil, "DELETE", path, nil)
}

// doAPIRequest creates and perform a standard HTTP request.
func doAPIRequest(v interface{}, method, path string, body interface{}) error {

	// the request URL includes the APIURL + APIVersion + path + user_slug + auth_token
	reqPath := APIURL + "/" + APIVersion + path + "?user_slug=" + UserSlug + "&auth_token=" + AuthToken

	req, err := NewRequest(method, reqPath, body, nil)
	if err != nil {
		return err
	}

	return Do(req, v)
}

// DoRawRequest creates and perform a standard HTTP request, allowing for the
// addition of custom headers
func DoRawRequest(v interface{}, method, path string, body interface{}, headers map[string]string) error {

	req, err := NewRequest(method, path, body, headers)
	if err != nil {
		return err
	}

	return Do(req, v)
}

// NewRequest creates an HTTP request for the nanobox API, but does not perform
// it.
func NewRequest(method, path string, body interface{}, headers map[string]string) (*http.Request, error) {

	var rbody io.Reader

	//
	switch t := body.(type) {
	case string:
		rbody = bytes.NewBufferString(t)
	case io.Reader:
		rbody = t
	default:
		rbody = nil
	}

	// an HTTP request
	req, err := http.NewRequest(method, path, rbody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", DefaultContentType)
	req.Header.Set("Content-Type", DefaultContentType)

	// add additional headers
	if headers != nil {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
	}

	return req, nil
}

// Do performs an http.NewRequest
func Do(req *http.Request, v interface{}) error {

	// debugging
	if Debug {
		dump, err := httputil.DumpRequestOut(req, true)
		if err != nil {
			return err
		}

		fmt.Println(`
Request:
--------------------------------------------------------------------------------
` + string(dump))
	}

	//
	res, err := HTTPClient.Do(req)
	if err != nil {
		return err
	}

	// debugging
	if Debug {
		dump, err := httputil.DumpResponse(res, true)
		if err != nil {
			return err
		}

		fmt.Println(`
Response:
--------------------------------------------------------------------------------
` + string(dump))
	}

	// read the body
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	// check the response
	if err = checkResponse(res, b); err != nil {
		return err
	}

	// unmarshal response into the provided container. If no container was given
	// it's mostly likely a raw request where the body isn't needed, and therfore
	// this step can be skipped.
	if v != nil {
		if err := json.Unmarshal(b, &v); err != nil {
			return err
		}
	}

	return err
}

// checkResponse determins if the response is !20* and return a custom api error
func checkResponse(res *http.Response, b []byte) error {

	if res.StatusCode/100 != 2 {

		apiError := APIError{
			error: errors.New(string(b)),
			Code:  res.StatusCode,
		}

		if err := json.Unmarshal(b, &apiError); err != nil {
			return err
		}

		return apiError
	}

	return nil
}

// helpers

// Bool is a helper that allocates a new bool value to store v and returns a
// pointer to it
func Bool(v bool) *bool {
	p := new(bool)
	*p = v
	return p
}

// Int is a helper that allocates a new int value to store v and returns a
// pointer to it
func Int(v int) *int {
	p := new(int)
	*p = v
	return p
}

// String is a helper that allocates a new string value to store v and returns a
// pointer to it
func String(v string) *string {
	p := new(string)
	*p = v
	return p
}
