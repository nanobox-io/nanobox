package models

// import (
// 	"testing"
// )

// func TestAuthSave(t *testing.T) {
// 	// clear the registry table when we're finished
// 	defer truncate("registry")

// 	auth := Auth{
// 		Key: "123",
// 	}

// 	err := auth.Save()
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	// fetch the auth
// 	auth2 := Auth{}

// 	if err = get("registry", "auth", &auth2); err != nil {
// 		t.Errorf("failed to fetch auth: %s", err.Error())
// 	}

// 	if auth2.Key != "123" {
// 		t.Errorf("auth doesn't match")
// 	}
// }

// func TestAuthDelete(t *testing.T) {
// 	// clear the registry table when we're finished
// 	defer truncate("registry")

// 	auth := Auth{
// 		Key: "123",
// 	}

// 	if err := auth.Save(); err != nil {
// 		t.Error(err)
// 	}

// 	if err := auth.Delete(); err != nil {
// 		t.Error(err)
// 	}

// 	// make sure the auth is gone
// 	keys, err := keys("registry")
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	if len(keys) > 0 {
// 		t.Errorf("auth was not deleted")
// 	}
// }

// func TestLoadAuth(t *testing.T) {
// 	// clear the registry table when we're finished
// 	defer truncate("registry")

// 	auth := Auth{
// 		Key: "123",
// 	}

// 	if err := auth.Save(); err != nil {
// 		t.Error(err)
// 	}

// 	auth2, err := LoadAuth()
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	if auth2.Key != "123" {
// 		t.Errorf("did not load the correct auth")
// 	}
// }
