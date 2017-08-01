package odin

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/loganmac/nanobox/models"
)

// Client holds the state and is responsible for making requests to Odin
type Client struct {
	HTTP            httpClient // the configured HTTP client
	URL             string     // which URL to send to
	BonesaltURL     string     // Bonesalt target URL
	DevURL          string     // Dev target URL
	SimURL          string     // Sim target URL
	NanoboxURL      string     // Nanobox (default) URL
	NanoboxUsername string     // username to connect to odin
	NanoboxPassword string     // password to connect to odin
	// inject an AuthRepo to retrieve stored auth tokens for various endpoints,
	// instead of reaching out to some global database
	AuthRepo AuthRepo
}

// httpClient is an interface that let's us mock out the
// more granular details of connecting to Odin and just
// return specific responses
type httpClient interface {
	do(method, url string, payload interface{}) (*odinResponse, error)
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

// SetTarget takes in a string and sets the client's current request URL to the
// one that corresponds with that slug, defaulting to the Nanobox URL
func (c *Client) SetTarget(s string) {
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
	// Try checking cached token
	var authToken string
	if err := c.AuthRepo.Get("auths", c.URL, &authToken); err != nil {
		return "", err
	}
	// return it if we have a cached one
	if authToken != "" {
		return authToken, nil
	}
	// Otherwise, we must not have a cached token, so let's get one

	// URLEncode the password
	var params url.Values
	params.Set("password", c.NanoboxPassword)
	// Construct the URL
	reqURL := fmt.Sprintf("%s/users/%s/auth_token?%s", c.URL, c.NanoboxUsername, params.Encode())
	// make request
	res, err := c.HTTP.do("GET", reqURL, nil)
	if err != nil {
		return "", nil
	}
	// cache authToken
	if err := c.AuthRepo.Put("auths", c.URL, res.AuthToken); err != nil {
		return "", err
	}
	return res.AuthToken, nil
}

// App takes in a slug like "appID" or "teamID/appID",
// and returns the app details from Odin
func (c *Client) App(slug string) (models.App, error) {
	app := models.App{}
	// construct parameters
	var params url.Values
	// split up teamID if it exists and add it to the params
	appendTeamContext(&slug, &params)
	// grab the auth token
	authToken, err := c.Auth()
	if err != nil {
		return app, err
	}
	params.Set("auth_token", authToken)
	// construct the URL
	reqURL := fmt.Sprintf("%s/apps/%s?%s", c.URL, slug, params.Encode())
	// make request
	res, err := c.HTTP.do("GET", reqURL, nil)
	if err != nil {
		return app, nil
	}
	_ = res // until we find out
	// TODO: find out the specific return of the app request to odin
	// so that we can marshal the response here and set it.
	return app, nil
}

// appendTeamContext is a helper to check the input for
// teamID/appID pattern, split it up, and append the teamID to the
// parameters as "ci"
func appendTeamContext(input *string, params *url.Values) {
	// check for the "teamname/appname" pattern
	if strings.Contains(*input, "/") {
		// split into [teamID, appID]
		teamAppIDs := strings.Split(*input, "/")
		// set the team name in the params under "ci"
		// for "Context ID"
		params.Set("ci", teamAppIDs[0])
		input = &teamAppIDs[1]
	}
}
