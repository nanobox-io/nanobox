package odin

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// Client holds the state and is responsible for making requests to Odin
type Client struct {
	HTTP            http.Client // the configured HTTP client
	URL             string      // which URL to send to
	BonesaltURL     string      // Bonesalt target URL
	DevURL          string      // Dev target URL
	SimURL          string      // Sim target URL
	NanoboxURL      string      // Nanobox (default) URL
	NanoboxUsername string      // username to connect to odin
	NanoboxPassword string      // password to connect to odin
	// inject an AuthRepo to retrieve stored auth tokens for various endpoints,
	// instead of reaching out to some global database
	AuthRepo AuthRepo
	// inject a logger so that we know what's going on, and have a logger that we
	// can configure one time, for things like setting the loglevel so you don't
	// have to inspect a debug flag in our application logic
	Logger logger
}

// AuthRepo is any key-value store, really, used here
// to cache the authToken for an endpoint
type AuthRepo interface {
	Get(bucket, id string, v interface{}) error
	Put(bucket, id string, v interface{}) error
}

// logger is an interface that can be implemented by any logger we inject,
// allowing us to test logging or swap out for a different one without touching
// code here
type logger interface {
	Debug(msg string, ctx ...interface{})
	Info(msg string, ctx ...interface{})
	Warn(msg string, ctx ...interface{})
	Error(msg string, ctx ...interface{})
}

// odinResponse is just a json decoding struct for any values Odin could respond with.
type odinResponse struct {
	Token     string `json:"token"`
	AuthToken string `json:"authentication_token"`
	URL       string `json:"url"`
	Protocol  string `json:"protocol"`
}

// SetEndpoint takes in a string and sets the client's current request URL to the
// one that corresponds with that slug, defaulting to the Nanobox URL
func (c *Client) SetEndpoint(s string) {
	endpoint := strings.ToUpper(s)
	switch endpoint {
	case "BONESALT":
		c.URL = c.BonesaltURL
	case "DEV":
		c.URL = c.DevURL
	case "SIM":
		c.URL = c.SimURL
	default:
		c.URL = c.NanoboxURL
	}
}

// Auth makes a request to the endpoint with the client's username and password,
// caches the returned authentication token, and returns it
func (c *Client) Auth() (string, error) {
	// URLEncode the password
	var params url.Values
	params.Set("password", c.NanoboxPassword)
	// Construct the URL
	URL := fmt.Sprintf("%s/users/%s/auth_token", c.URL, c.NanoboxUsername)
	// make request
	res, err := c.do("GET", URL, params, nil)
	if err != nil {
		return "", nil
	}
	// cache authToken
	c.AuthRepo.Put("auths", c.URL, res.AuthToken)
	return res.AuthToken, nil
}

// do wraps Client.HTTP.do with things like logging, json parsing, authentication, and error handling
func (c *Client) do(method, url string, params url.Values, payload interface{}) (*odinResponse, error) {
	var reqBody *bytes.Buffer
	// marshal the payload to json, if there is one
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
	// get auth key for endpoint
	var authToken string
	c.AuthRepo.Get("auths", c.URL, &authToken)
	if authToken == "" {
		// try authenticating if there isn't one
		authToken, err = c.Auth()
		if err != nil {
			return nil, err
		}
	}
	// if we are debugging, log request
	c.Logger.Debug("Odin Request", "method", req.Method, "url", req.URL, "proto", req.Proto)

	// call upon the might of the Alfadir
	res, err := c.HTTP.Do(req)
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
	c.Logger.Debug("Odin Response", "body", resBody, "statusCode", res.StatusCode, "method", req.Method,
		"url", req.URL, "proto", req.Proto, "content-length", res.Header.Get("Content-Length"))

	// handle error responses
	if res.StatusCode < 200 || res.StatusCode > 299 {
		c.Logger.Error("bad response from Odin", "body", resBody, "statusCode", res.StatusCode, "method", req.Method,
			"url", req.URL, "proto", req.Proto, "content-length", res.Header.Get("Content-Length"))
		return nil, errors.New("received a bad response from odin")
	}

	// decode the response json
	var odinResp odinResponse
	if err := json.Unmarshal(resBody, &odinResp); err != nil {
		// sometimes it's difficult for mere mortals to comprehend the Alfadir's
		// mighty tambre, so we just log that out
		c.Logger.Error("error parsing response from Odin", "error", err)
		return nil, fmt.Errorf("could not parse response: %v", err)
	}
	return &odinResp, nil
}
