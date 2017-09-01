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
	"os"
	"strings"

	"github.com/jcelliott/lumber"

	"github.com/nanobox-io/nanobox/commands/registry"
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util"
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
	apiKey   string
)

type (
	evar struct {
		ID    string `json:"id"`
		Key   string `json:"title"`
		Value string `json:"value"`
	}
)

// sets the odin endpoint
func SetEndpoint(stage string) {
	endpoint = stage
}

// Auth authenticates the user with odin.
func Auth(username, password string) (string, error) {

	loginInfo := struct {
		Slug     string `json:"slug"`
		Password string `json:"password"`
	}{username, password}

	resBody := map[string]string{}

	if err := doRequest("GET", "user_auth_token", nil, loginInfo, &resBody); err != nil {
		return "", err
	}

	return resBody["authentication_token"], nil
}

// App ...
func App(slug string) (models.App, error) {
	app := models.App{}
	var params url.Values
	if strings.Contains(slug, "/") {
		appNameParts := strings.Split(slug, "/")
		if len(appNameParts) == 2 {
			params = url.Values{}
			params.Set("ci", appNameParts[0])
			slug = appNameParts[1]
		}

	}

	return app, doRequest("GET", "apps/"+slug, params, nil, &app)
}

// Deploy ...
func Deploy(appID, id, boxfile, message string) error {

	//
	body := map[string]map[string]string{
		"deploy": {
			"boxfile_content": boxfile,
			"build_id":        id,
			"commit_message":  message,
		},
	}

	var params url.Values

	if strings.Contains(appID, "/") {
		appNameParts := strings.Split(appID, "/")
		if len(appNameParts) == 2 {
			params = url.Values{}
			params.Set("ci", appNameParts[0])
			appID = appNameParts[1]
		}

	}

	return doRequest("POST", fmt.Sprintf("apps/%s/deploys", appID), params, body, nil)
}

func ListEvars(appID string) ([]evar, error) {
	evars := []evar{}

	var params url.Values
	if strings.Contains(appID, "/") {
		appNameParts := strings.Split(appID, "/")
		if len(appNameParts) == 2 {
			params = url.Values{}
			params.Set("ci", appNameParts[0])
			appID = appNameParts[1]
		}

	}

	return evars, doRequest("GET", fmt.Sprintf("apps/%s/evars", appID), params, nil, &evars)
}

func AddEvar(appID, key, val string) error {
	body := map[string]map[string]string{
		"evar": {
			"title": key,
			"value": val,
		},
	}

	var params url.Values

	if strings.Contains(appID, "/") {
		appNameParts := strings.Split(appID, "/")
		if len(appNameParts) == 2 {
			params = url.Values{}
			params.Set("ci", appNameParts[0])
			appID = appNameParts[1]
		}

	}

	return doRequest("POST", fmt.Sprintf("apps/%s/evars", appID), params, body, nil)
}

func RemoveEvar(appID, id string) error {
	var params url.Values

	if strings.Contains(appID, "/") {
		appNameParts := strings.Split(appID, "/")
		if len(appNameParts) == 2 {
			params = url.Values{}
			params.Set("ci", appNameParts[0])
			appID = appNameParts[1]
		}

	}

	return doRequest("DELETE", fmt.Sprintf("apps/%s/evars/%s", appID, id), params, nil, nil)
}

// EstablishTunnel requests a tunnel from odin.
func EstablishTunnel(tunCfg models.TunnelConfig) (models.TunnelInfo, error) {
	r := models.TunnelInfo{Port: tunCfg.DestPort}

	params := url.Values{}
	if strings.Contains(tunCfg.AppName, "/") {
		appNameParts := strings.Split(tunCfg.AppName, "/")
		if len(appNameParts) == 2 {
			params = url.Values{}
			// ci is contextID for teams `team/app`
			params.Set("ci", appNameParts[0])
			tunCfg.AppName = appNameParts[1]
		}
	}

	err := doRequest("GET", fmt.Sprintf("apps/%s/tunnels/%s", tunCfg.AppName, tunCfg.Component), params, r, &r)
	if err != nil {
		if strings.Contains(err.Error(), "valid port") {
			if err2, ok := err.(util.Err); ok {
				err2.Suggest = "Specify a valid destination port with `-p source:destination`"
				err2.Message = err.Error()
				return r, err2
			}
		}
		if strings.Contains(err.Error(), "no default") {
			if err2, ok := err.(util.Err); ok {
				err2.Suggest = "Specify a destination port with `-p source:destination`"
				err2.Message = err.Error()
				return r, err2
			}
		}
		if strings.Contains(err.Error(), "Not Found") {
			if err2, ok := err.(util.Err); ok {
				err2.Message = fmt.Sprintf("%s - component '%s' not found on app '%s'", err.Error(), tunCfg.Component, tunCfg.AppName)
				return r, err2
			}
			err = fmt.Errorf("%s - component '%s' not found on app '%s'", err.Error(), tunCfg.Component, tunCfg.AppName)
		}
	}

	return r, err
}

// EstablishConsole ...
// protocol ssh/docker
func EstablishConsole(appID, id string) (string, string, string, error) {
	// use a default user
	params := url.Values{}
	params.Set("user", "gonano")
	if registry.GetString("console_user") != "" {
		params.Set("user", registry.GetString("console_user"))
	}

	if strings.Contains(appID, "/") {
		appNameParts := strings.Split(appID, "/")
		if len(appNameParts) == 2 {
			params.Set("ci", appNameParts[0])
			appID = appNameParts[1]
		}

	}

	r := map[string]string{}
	err := doRequest("GET", fmt.Sprintf("apps/%s/consoles/%s", appID, id), params, nil, &r)

	return r["token"], r["url"], r["protocol"], err
}

// GetWarehouse ...
func GetWarehouse(appID string) (string, string, error) {
	r := map[string]string{}

	var params url.Values
	if strings.Contains(appID, "/") {
		appNameParts := strings.Split(appID, "/")
		if len(appNameParts) == 2 {
			params = url.Values{}
			params.Set("ci", appNameParts[0])
			appID = appNameParts[1]
		}

	}

	err := doRequest("GET", fmt.Sprintf("apps/%s/services/warehouse", appID), params, nil, &r)

	return r["token"], r["url"], err
}

// GetComponent gets the token and url for a given component
func GetComponent(appID, component string) (string, string, error) {
	r := map[string]string{}

	var params url.Values
	if strings.Contains(appID, "/") {
		appNameParts := strings.Split(appID, "/")
		if len(appNameParts) == 2 {
			params = url.Values{}
			params.Set("ci", appNameParts[0])
			appID = appNameParts[1]
		}

	}

	err := doRequest("GET", fmt.Sprintf("apps/%s/services/%s", appID, component), params, nil, &r)
	// r = map["uid":"pusher1" "name":"message bus" "slug":"pusher" "url":"x.x.x.x" "token":"secret" "mode":"simple" "ip":"x.x.x.x" "id":"d371c59b-fbbc-4248-8f3b-7f694a62b5e1"]

	return r["token"], r["url"], err
}

func GetPreviousBuild(appID string) (string, error) {
	r := []map[string]string{}

	var params url.Values
	if strings.Contains(appID, "/") {
		appNameParts := strings.Split(appID, "/")
		if len(appNameParts) == 2 {
			params = url.Values{}
			params.Set("ci", appNameParts[0])
			appID = appNameParts[1]
		}

	}

	err := doRequest("GET", fmt.Sprintf("apps/%s/deploys", appID), params, nil, &r)
	if err != nil {
		return "", err
	}

	if len(r) > 0 {
		return r[0]["build_id"], nil
	}

	return "", nil
}

func SubmitEvent(action, message, app string, meta map[string]interface{}) error {
	params := url.Values{}
	params.Set("api_key", apiKey)

	request := struct {
		Action  string                 `json:"action"`
		App     string                 `json:"eventable_id,omitempty"`
		Meta    map[string]interface{} `json:"meta"`
		Message string                 `json:"message"`
	}{
		Action:  action,
		App:     app,
		Meta:    meta,
		Message: message,
	}

	err := doRequest("POST", "events", params, map[string]interface{}{"event": request}, nil)
	if err != nil {
		return err
	}

	return nil
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

	// if they have not logged in but the user name and password are both set
	// use attempt to authenticate
	if auth.Key == "" &&
		os.Getenv("NANOBOX_PASSWORD") != "" &&
		os.Getenv("NANOBOX_USERNAME") != "" &&
		path != fmt.Sprintf("users/%s/auth_token", os.Getenv("NANOBOX_USERNAME")) {

		auth.Key, _ = Auth(os.Getenv("NANOBOX_USERNAME"), os.Getenv("NANOBOX_PASSWORD"))
	}

	if params == nil {
		params = url.Values{}
	}
	if auth.Key != "" {
		params.Set("auth_token", auth.Key)
	}

	// fetch the correct url from the endpoint
	url := odinURL()

	//
	lumber.Debug("%s%s?%s\n", url, path, params.Encode())
	req, err := http.NewRequest(method, fmt.Sprintf("%s%s?%s", url, path, params.Encode()), rbodyReader)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	lumber.Trace("REQ: %s %s %s", req.Method, req.URL, req.Proto)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	lumber.Debug("RES: %d %s %s %s (%s)", res.StatusCode, req.Method, req.URL, req.Proto, res.Header.Get("Content-Length"))

	// print the body even if status is not 2XX
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if res.StatusCode == 401 {
		err = util.ErrorfQuiet("[USER] Unauthorized (%s)", b)
		if err != nil {
			if err2, ok := err.(util.Err); ok {
				err2.Suggest = "It appears you are not logged in, run `nanobox login` and try again"
				return err2
			}
		}
		return err
	}

	if res.StatusCode == 404 {
		err = util.ErrorfQuiet("[USER] Not Found (%s)", b)
		if err != nil {
			if err2, ok := err.(util.Err); ok {
				err2.Suggest = "It appears that app or component wasn't found, check the team/app name and try again"
				return err2
			}
		}
		return err
	}

	// if it is a 400 but not
	if res.StatusCode >= 400 && res.StatusCode < 500 {
		rb := map[string]string{}
		err = json.Unmarshal(b, &rb)
		if err != nil {
			return util.ErrorfQuiet("[NANOBOX] %s", b)
		}

		errorMessage, ok := rb["error"]
		if !ok {
			return util.ErrorfQuiet("[ODIN] %s", b)
		}
		return util.ErrorfQuiet("[ODIN] %s", errorMessage)
	}

	if res.StatusCode == 500 {
		return util.ErrorfQuiet("[ODIN] Internal Server Error (%s)", b)
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return util.ErrorfQuiet("[ODIN] bad exit response(%d %s %s %s (%s) %s)", res.StatusCode, req.Method, req.URL, req.Proto, res.Header.Get("Content-Length"), b)
	}

	if responseBody != nil {
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
