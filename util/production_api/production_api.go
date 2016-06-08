package production_api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util/data"
	"github.com/nanobox-io/nanobox/util/config"
)

func Auth(username, password string) (string, error) {
	reqBody := map[string]string{
		"slug":     username,
		"password": password,
	}
	resBody := map[string]string{}
	err := doRequest("GET", "users/"+username+"/auth_token", reqBody, &resBody)
	if err != nil {
		return "", err
	}

	return resBody["authentication_token"], nil
}

func App(slug string) (models.App, error) {
	app := models.App{}
	return app, doRequest("GET", "apps/"+slug, nil, &app)
}

func Deploy(appId, id, boxfile, message string) error {
	body := map[string]string{
		"boxfile_content": boxfile,
		"build_id":        id,
		"commit_message":  message,
	}
	return doRequest("POST", fmt.Sprintf("/apps/%s/deploys", appId), body, nil)
}

func EstablishTunnel(appId, id string) (string, string, string, error) {
	// do somethign else here
	r := map[string]string{}
	err := doRequest("Get", fmt.Sprintf("/apps/%s/tunnels/%s", appId, id), nil, &r)
	return r["token"], r["url"], r["container"], err
}

func EstablishConsole(appId, id string) (string, string, string, error) {
	// do somethign else here
	r := map[string]string{}
	err := doRequest("Get", fmt.Sprintf("/apps/%s/consoles/%s", appId, id), nil, &r)
	return r["token"], r["url"], r["container"], err
}

func GetWarehouse(appId string) (string, string, error) {
	// do something else here
	r := map[string]string{}
	err := doRequest("Get", fmt.Sprintf("/apps/%s/warehouse", appId), nil, &r)
	return r["token"], r["url"], err
}

func doRequest(method, path string, requestBody, responseBody interface{}) error {
	var rbodyReader io.Reader
	if requestBody != nil {
		jsonBytes, err := json.Marshal(requestBody)
		if err != nil {
			return err
		}
		rbodyReader = bytes.NewBuffer(jsonBytes)
	}

	host := config.Viper().GetString("production_host")
	auth := models.Auth{}
	data.Get("global", "user", &auth)

	req, err := http.NewRequest(method, fmt.Sprintf("https://%s/%s?auth_token=%s", host, path, auth.Key), rbodyReader)
	if err != nil {
		return err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if responseBody != nil {
		b, err := ioutil.ReadAll(res.Body)
		err = json.Unmarshal(b, responseBody)
		if err != nil {
			return err
		}
	}
	return nil
}
