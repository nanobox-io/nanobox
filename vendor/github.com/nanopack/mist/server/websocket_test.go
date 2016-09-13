package server_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/nanopack/mist/auth"
	"github.com/nanopack/mist/server"
)

// TestWSStart tests to ensure a server will start
func TestWSStart(t *testing.T) {
	fmt.Println("Starting WS test...")

	// ensure authentication is disabled
	auth.DefaultAuth = nil

	go func() {
		if err := server.Start([]string{"ws://127.0.0.1:8888"}, ""); err != nil {
			t.Fatalf("Unexpected error - %v", err.Error())
		}
	}()
	<-time.After(time.Second)
}

// TestWSSStart tests to ensure a server will start
func TestWSSStart(t *testing.T) {
	fmt.Println("Starting WSS test...")

	// ensure authentication is disabled
	auth.DefaultAuth = nil

	go func() {
		if err := server.Start([]string{"wss://127.0.0.1:8988"}, ""); err != nil {
			t.Fatalf("Unexpected error - %v", err.Error())
		}
	}()
	<-time.After(time.Second)
}
