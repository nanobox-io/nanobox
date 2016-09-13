package server_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/nanopack/mist/auth"
	"github.com/nanopack/mist/server"
)

// TestHTTPStart tests to ensure a server will start
func TestHTTPStart(t *testing.T) {
	fmt.Println("Starting HTTP test...")

	// ensure authentication is disabled
	auth.DefaultAuth = nil
	go func() {
		if err := server.Start([]string{"http://127.0.0.1:8080"}, ""); err != nil {
			t.Fatalf("Unexpected error - %v", err.Error())
		}
	}()
	<-time.After(time.Second)
}
