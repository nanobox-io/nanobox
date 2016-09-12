// Package odin ...
package odin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/jcelliott/lumber"

	"github.com/nanobox-io/nanobox/models"
)

const (
	NANOBOX  = "https://api.nanobox.io/v1/"
	BONESALT = "https://api.bonesalt.com/v1/"
	DEV      = "http://api.nanobox.dev:8080/v1/"
	SIM      = "http://api.nanobox.sim/v1/"
)

var (
	// set the default endpoint to nanobox
	endpoint = "nanobox"
)

// sets the odin endpoint
func SetEndpoint(stage string) {
	endpoint = stage
}

// Auth ...
func Auth(username, password string) (string, error) {

	//
	params := url.Values{}
	params.Set("password", password)

	//
	resBody := map[string]string{}

	//
	if err := doRequest("GET", fmt.Sprintf("users/%s/auth_token", username), params, nil, &resBody); err != nil {
		return "", err
	}

	return resBody["authentication_token"], nil
}

// App ...
func App(slug string) (models.App, error) {
	app := models.App{}

	return app, doRequest("GET", "apps/"+slug, nil, nil, &app)
}

// Deploy ...
func Deploy(appID, id, boxfile, message string) error {

	//
	body := map[string]map[string]string{
		"deploy": map[string]string{
			"boxfile_content": boxfile,
			"build_id":        id,
			"commit_message":  message,
		},
	}

	return doRequest("POST", fmt.Sprintf("apps/%s/deploys", appID), nil, body, nil)
}

// EstablishTunnel ...
func EstablishTunnel(appID, id string) (string, string, int, error) {
	r := struct {
		Token string `json:"token"`
		Url   string `json:"url"`
		Port  int    `json:"port"`
	}{}

	err := doRequest("GET", fmt.Sprintf("apps/%s/tunnels/%s", appID, id), nil, nil, &r)

	return r.Token, r.Url, r.Port, err
}

// EstablishConsole ...
// protocol ssh/docker
func EstablishConsole(appID, id string) (string, string, string, error) {
	r := map[string]string{}
	err := doRequest("GET", fmt.Sprintf("apps/%s/consoles/%s", appID, id), nil, nil, &r)

	return r["token"], r["url"], r["protocol"], err
}

// GetWarehouse ...
func GetWarehouse(appID string) (string, string, error) {
	r := map[string]string{}
	err := doRequest("GET", fmt.Sprintf("apps/%s/services/warehouse", appID), nil, nil, &r)

	return r["token"], r["url"], err
}

func GetPreviousBuild(appID string) (string, error) {
	r := []map[string]string{}
	err := doRequest("GET", fmt.Sprintf("apps/%s/deploys", appID), nil, nil, &r)
	if err != nil {
		return "", err
	}

	if len(r) > 0 {
		return r[0]["build_id"], nil
	}

	return "", nil
}

// doRequest ...
func doRequest(method, path string, params url.Values, requestBody, responseBody interface{}) error {

	var rbodyReader io.Reader

	//
	if requestBody != nil {
		jsonBytes, err := json.Marshal(requestBody)
		if err != nil {
			return err
		}
		rbodyReader = bytes.NewBuffer(jsonBytes)
	}

	auth, _ := models.LoadAuthByEndpoint(endpoint)

	if params == nil {
		params = url.Values{}
	}
	params.Set("auth_token", auth.Key)

	// fetch the correct url from the endpoint
	url := odinURL()

	//
	lumber.Debug("%s%s?%s\n", url, path, params.Encode())
	req, err := http.NewRequest(method, fmt.Sprintf("%s%s?%s", url, path, params.Encode()), rbodyReader)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	lumber.Debug("req: %+v\n", req)

	//
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	lumber.Debug("res: %+v\n", res)

	if res.StatusCode == 401 {
		return fmt.Errorf("Unauthorized")
	}

	if res.StatusCode == 404 {
		return fmt.Errorf("Not Found")
	}

	if res.StatusCode == 500 {
		return fmt.Errorf("Internal Server Error")
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return fmt.Errorf("bad exit response(%+v)", res)
	}

	//
	if responseBody != nil {
		b, err := ioutil.ReadAll(res.Body)
		lumber.Debug("response body: '%s'\n", b)
		err = json.Unmarshal(b, responseBody)
		if err != nil {
			return err
		}
	}

	return nil
}

func odinURL() string {
	switch endpoint {
	case "bonesalt":
		return BONESALT
	case "dev":
		return DEV
	case "sim":
		return SIM
	default:
		return NANOBOX
	}
}
