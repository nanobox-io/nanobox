package models

import (
	"testing"
)

type data struct {
	Name   string
	Number int
}

func init() {
	// initialize the db in the .nanobox directory
	DB = "/tmp/nanobox-test.db"
}

func TestPut(t *testing.T) {
	d := data{
		Name:   "mickey",
		Number: 1234,
	}

	err := put("user", "1", d)
	if err != nil {
		t.Errorf("unable to put data in bucket %+v", err)
	}
}

func TestGet(t *testing.T) {
	d := data{}

	err := get("user", "1", &d)
	if err != nil {
		t.Errorf("error getting data %+v", err)
	}

	if d.Name != "mickey" || d.Number != 1234 {
		t.Errorf("retrieved data does not match %+v", d)
	}
}

func TestDelete(t *testing.T) {
	// clear the users table when we're finished
	defer truncate("users")

	err := destroy("user", "1")
	if err != nil {
		t.Errorf("unable to delete %+v", err)
	}

	d := data{}
	err = get("user", "1", &d)
	if err == nil {
		t.Errorf("removed data should not have been retrievable")
	}
}

func TestKeys(t *testing.T) {
	// clear the users table when we're finished
	defer truncate("users")

	mickey := data{Name: "Mickey"}
	minnie := data{Name: "Minnie"}
	donald := data{Name: "Donald"}

	put("user", "1", mickey)
	put("user", "2", minnie)
	put("user", "3", donald)

	users, err := keys("user")
	if err != nil {
		t.Errorf("failed to list keys for 'user' bucket")
	}

	if len(users) != 3 {
		t.Errorf("failed to list all keys for 'user' bucket")
	}
}

func TestTruncate(t *testing.T) {
	mickey := data{Name: "Mickey"}
	minnie := data{Name: "Minnie"}

	put("user", "1", mickey)
	put("user", "2", minnie)

	if err := truncate("users"); err != nil {
		t.Errorf("failed to truncate users: %s", err.Error())
	}

	users, err := keys("users")
	if err != nil {
		t.Errorf("failed to list keys for 'user' bucket: %s", err.Error())
	}

	if len(users) != 0 {
		t.Errorf("'users' bucket was not truncated")
	}
}
