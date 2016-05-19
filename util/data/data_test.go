package data_test

import "testing"
import "github.com/nanobox-io/nanobox/util/data"

type Data struct {
	Name   string
	Number int
	Cows   bool
}

func TestPut(t *testing.T) {
	d := Data{
		"lyon",
		1234,
		true,
	}
	err := data.Put("user", "1", d)
	if err != nil {
		t.Errorf("unable to put data in bucket %+v", err)
	}
}

func TestGet(t *testing.T) {
	d := Data{}
	err := data.Get("user", "1", &d)
	if err != nil {
		t.Errorf("error getting data %+v", err)
	}
	if d.Name != "lyon" || d.Number != 1234 || !d.Cows {
		t.Errorf("retrieved data does not match %+v", d)
	}
}

func TestDelete(t *testing.T) {
	err := data.Delete("user", "1")
	if err != nil {
		t.Errorf("unable to delete %+v", err)
	}
	d := Data{}
	err = data.Get("user", "1", &d)
	if err == nil {
		t.Errorf("removed data should not have been retrievable")
	}

}
