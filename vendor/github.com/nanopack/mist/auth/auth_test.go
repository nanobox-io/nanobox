package auth_test

import (
	"net/url"
	"os"
	"testing"

	_ "github.com/lib/pq"

	mistAuth "github.com/nanopack/mist/auth"
)

var (
	testToken = "token"
	testTag1  = "onefish"
	testTag2  = "twofish"
	testTag3  = "redfish"
	testTag4  = "bluefish"
)

// TestStart tests the auth start process
func TestStart(t *testing.T) {

	//
	if err := mistAuth.Start("memory://"); err != nil {
		t.Fatalf("Unexpected error!")
	}

	// DefaultAuth is set inside of an auth start and should not be nil once started
	if mistAuth.DefaultAuth == nil {
		t.Fatalf("Unexpected nil DefaultAuth!")
	}
}

// TestMemory tests the memory authenticator
func TestMemory(t *testing.T) {

	//
	url, err := url.Parse("memory://")
	if err != nil {
		t.Fatalf(err.Error())
	}

	// create a new memory authenticator
	mem, err := mistAuth.NewMemory(url)
	if err != nil {
		t.Fatalf(err.Error())
	}

	//
	testAuth(mem, t)
}

// TestScribble tests the scribble authenticator
func TestScribble(t *testing.T) {

	// attempt to remove the db from any previous tests
	if err := os.RemoveAll("/tmp/scribble"); err != nil {
		t.Fatalf(err.Error())
	}

	//
	url, err := url.Parse("scribble://?db=/tmp/scribble")
	if err != nil {
		t.Fatalf(err.Error())
	}

	//
	scribble, err := mistAuth.NewScribble(url)
	if err != nil {
		t.Fatalf(err.Error())
	}

	//
	testAuth(scribble, t)
}

// TestPostgres tests the postgres authenticator (requires running postgres server)
// func TestPostgres(t *testing.T) {
//
// 	//
// 	url, err := url.Parse("postgres://postgres@127.0.0.1:5432?db=postgres")
// 	if err != nil {
// 		t.Fatalf(err.Error())
// 	}
//
// 	//
// 	pg, err := NewPostgres(url)
// 	if err != nil {
// 		t.Fatalf(err.Error())
// 	}
//
// 	//
// 	if _, err := pg.(postgresql).exec("TRUNCATE tokens, tags"); err != nil {
// 		t.Fatalf(err.Error())
// 	}
//
// 	//
// 	testAuth(pg, t)
// }

// testAuth tests to ensure all authenticator methods are working as intended
func testAuth(auth mistAuth.Authenticator, t *testing.T) {

	// no token should exist yet
	tags, err := auth.GetTagsForToken(testToken)
	if err == nil {
		t.Fatalf("Unexpected token!")
	}
	if len(tags) != 0 {
		t.Fatalf("Unexpected tags!")
	}

	// add a new token
	if err := auth.AddToken(testToken); err != nil {
		t.Fatalf(err.Error())
	}

	// add tags
	if err := auth.AddTags(testToken, []string{testTag1, testTag2}); err != nil {
		t.Fatalf(err.Error())
	}

	// add same tags; these should not get added
	if err := auth.AddTags(testToken, []string{testTag1, testTag2}); err != nil {
		t.Fatalf(err.Error())
	}

	// add same tags, different order; these should not get added
	if err := auth.AddTags(testToken, []string{testTag2, testTag1}); err != nil {
		t.Fatalf(err.Error())
	}

	// add more tags
	if err := auth.AddTags(testToken, []string{testTag3, testTag4}); err != nil {
		t.Fatalf(err.Error())
	}

	// this tests to ensure that same tags don't get added and we only get back
	// what we expect (in this case 4 unique tags)
	tags, err = auth.GetTagsForToken(testToken)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if len(tags) != 4 {
		t.Fatalf("Wrong number of tags. Expecting 4 received %v", len(tags))
	}

	// remove tags
	if err := auth.RemoveTags(testToken, []string{testTag1, testTag2}); err != nil {
		t.Fatalf(err.Error())
	}

	// remote token
	if err := auth.RemoveToken(testToken); err != nil {
		t.Fatalf(err.Error())
	}

	// same as the first test; the token should no longer exist
	tags, err = auth.GetTagsForToken(testToken)
	if err == nil {
		t.Fatalf("Unexpected token!")
	}
	if len(tags) != 0 {
		t.Fatalf("Unexpected tags")
	}
}
