package models

import (
	"testing"
)

func TestEnvSave(t *testing.T) {
	// clear the envs table when we're finished
	defer truncate("envs")

	env := Env{
		ID:        "123",
		Directory: "/foo/bar",
		Name:      "test-env",
	}

	err := env.Save()
	if err != nil {
		t.Error(err)
	}

	// fetch the env
	env2 := Env{}

	if err = get("envs", env.ID, &env2); err != nil {
		t.Errorf("failed to fetch env: %s", err.Error())
	}

	if env2.ID != "123" {
		t.Errorf("env doesn't match")
	}
}

func TestEnvDelete(t *testing.T) {
	// clear the envs table when we're finished
	defer truncate("envs")

	env := Env{
		ID:        "123",
		Directory: "/foo/bar",
		Name:      "test-env",
	}

	if err := env.Save(); err != nil {
		t.Error(err)
	}

	if err := env.Delete(); err != nil {
		t.Error(err)
	}

	// make sure the env is gone
	keys, err := keys("envs")
	if err != nil {
		t.Error(err)
	}

	if len(keys) > 0 {
		t.Errorf("env was not deleted")
	}
}

func TestFindEnvByID(t *testing.T) {
	// clear the envs table when we're finished
	defer truncate("envs")

	env := Env{
		ID:        "123",
		Directory: "/foo/bar",
		Name:      "test-env",
	}

	if err := env.Save(); err != nil {
		t.Error(err)
	}

	env2, err := FindEnvByID("123")
	if err != nil {
		t.Error(err)
	}

	if env2.ID != "123" {
		t.Errorf("did not load the correct env")
	}
}

func TestAllEnvs(t *testing.T) {
	// clear the envs table when we're finished
	defer truncate("envs")

	env1 := Env{ID: "1"}
	env2 := Env{ID: "2"}
	env3 := Env{ID: "3"}

	if err := env1.Save(); err != nil {
		t.Error(err)
	}
	if err := env2.Save(); err != nil {
		t.Error(err)
	}

	if err := env3.Save(); err != nil {
		t.Error(err)
	}

	envs, err := AllEnvs()
	if err != nil {
		t.Error(err)
	}

	if len(envs) != 3 {
		t.Errorf("did not load all envs")
	}
}
