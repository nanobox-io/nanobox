package odin

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

// HTTP handles all the fiddly pieces about communicating with odin
type HTTP struct {
	// inject a logger so that we know what's going on, and have a logger that we
	// can configure one time, for things like setting the loglevel so you don't
	// have to inspect a debug flag in our application logic
	Logger logger
}

// do wraps Client.HTTP.do with things like logging, json parsing, and error handling
func (h *HTTP) do(method, url string, payload interface{}) (*odinResponse, error) {
	// marshal the payload to json, if there is one
	var reqBody *bytes.Buffer
	if payload != nil {
		jsonBody, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}
	// construct the request to be sent to Odin
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// if we are debugging, log request
	h.Logger.Debug("Odin Request", "method", req.Method, "url", req.URL, "proto", req.Proto)

	// call upon the might of the Alfadir
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// read body into buffer
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	// if we are debugging, log response
	h.Logger.Debug("Odin Response", "body", resBody, "statusCode", res.StatusCode, "method", req.Method,
		"url", req.URL, "proto", req.Proto, "content-length", res.Header.Get("Content-Length"))

	// handle error responses
	if res.StatusCode < 200 || res.StatusCode > 299 {
		h.Logger.Error("bad response from Odin", "body", resBody, "statusCode", res.StatusCode, "method", req.Method,
			"url", req.URL, "proto", req.Proto, "content-length", res.Header.Get("Content-Length"))
		return nil, errors.New("received a bad response from odin")
	}

	// decode the response json
	var odinResp odinResponse
	if err := json.Unmarshal(resBody, &odinResp); err != nil {
		// sometimes it's difficult for mere mortals to comprehend the Alfadir's
		// mighty tambre, so we just log that out
		h.Logger.Error("error parsing response from Odin", "error", err)
		return nil, fmt.Errorf("could not parse response: %v", err)
	}
	return &odinResp, nil
}
