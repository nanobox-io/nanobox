// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

package util

//
import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
)

//
func S3RequestURL(path string) (string, error) {

	//
	res, err := http.DefaultClient.Get(path)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	//
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	preq := make(map[string]string)

	if err := json.Unmarshal(b, &preq); err != nil {
		return "", err
	}

	return preq["url"], nil
}

// S3Download
func S3Download(path string) (*http.Response, error) {
	return s3Request("GET", path, nil)
}

// S3Upload
func S3Upload(path string, body io.Reader) error {
	res, err := s3Request("PUT", path, body)
	defer res.Body.Close()

	return err
}

// s3Request
func s3Request(method, path string, body io.Reader) (*http.Response, error) {

	//
	s3req, err := http.NewRequest(method, path, body)
	if err != nil {
		return nil, err
	}

	//
	s3res, err := http.DefaultClient.Do(s3req)
	if err != nil {
		return nil, err
	}

	return s3res, nil
}
