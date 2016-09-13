package server_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/nanopack/mist/auth"
	"github.com/nanopack/mist/server"
)

// TestAuthStart tests an auth start process
func TestAuthStart(t *testing.T) {
	fmt.Println("Starting SERVER test...")

	// start an authenticator
	if err := auth.Start("memory://"); err != nil {
		t.Fatalf("Unexpected error - %v", err.Error())
	}

	// test for error if an auth is provided w/o a token
	go func() {
		if err := server.Start([]string{"tcp://127.0.0.1:1446"}, ""); err == nil {
			t.Fatalf("Expecting error - %v", err.Error())
		}
	}()
	<-time.After(time.Second)

	// test for successful start if token is provided
	go func() {
		if err := server.Start([]string{"tcp://127.0.0.1:1446"}, "TOKEN"); err != nil {
			t.Fatalf("Unexpected error - %v", err.Error())
		}
	}()
	<-time.After(time.Second)

	// test for error if authtoken does not match the token the server started with
	// if authtoken != "TOKEN" {
	// 	t.Fatalf("Incorrect token!")
	// }
}
