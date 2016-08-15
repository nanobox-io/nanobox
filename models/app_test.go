package models

import (
	"fmt"
	"testing"
)

func TestAppSave(t *testing.T) {
	// clear the envs table when we're finished
	defer truncate("123")

	app := App{
		EnvID: "123",
		ID:    "123_dev",
		Name:  "dev",
	}

	err := app.Save()
	if err != nil {
		t.Error(err)
	}

	// fetch the app
	app2 := App{}

	key := fmt.Sprintf("%s_%s", app.EnvID, app.Name)
	if err = get(app.EnvID, key, &app2); err != nil {
		t.Errorf("failed to fetch app: %s", err.Error())
	}

	if app2.EnvID != "123" {
		t.Errorf("app doesn't match")
	}
}

func TestAppDelete(t *testing.T) {
	// clear the envs table when we're finished
	defer truncate("123")

	app := App{
		EnvID: "123",
		ID:    "123_dev",
		Name:  "dev",
	}

	if err := app.Save(); err != nil {
		t.Error(err)
	}

	if err := app.Delete(); err != nil {
		t.Error(err)
	}

	// make sure the app is gone
	keys, err := keys(app.EnvID)
	if err != nil {
		t.Error(err)
	}

	if len(keys) > 0 {
		t.Errorf("app was not deleted")
	}
}

func TestFindAppBySlug(t *testing.T) {
	// clear the envs table when we're finished
	defer truncate("123")

	app := App{
		EnvID: "123",
		ID:    "123_dev",
		Name:  "dev",
	}

	if err := app.Save(); err != nil {
		t.Error(err)
	}

	app2, err := FindAppBySlug("123", "dev")
	if err != nil {
		t.Error(err)
	}

	if app2.EnvID != "123" || app2.Name != "dev" {
		t.Errorf("did not load the correct env")
	}
}

func TestAllAppsByEnv(t *testing.T) {
	// clear the envs table when we're finished
	defer truncate("1")
	defer truncate("2")

	env1 := App{EnvID: "1", ID: "1_dev", Name: "dev"}
	env2 := App{EnvID: "1", ID: "1_sim", Name: "sim"}
	env3 := App{EnvID: "2", ID: "2_dev", Name: "dev"}

	if err := env1.Save(); err != nil {
		t.Error(err)
	}
	if err := env2.Save(); err != nil {
		t.Error(err)
	}

	if err := env3.Save(); err != nil {
		t.Error(err)
	}

	envs, err := AllAppsByEnv("1")
	if err != nil {
		t.Error(err)
	}

	if len(envs) != 2 {
		t.Errorf("did not load all envs, got %d", len(envs))
	}
}
