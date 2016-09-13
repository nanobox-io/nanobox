package clients_test

import (
	"os"
	"testing"
	"time"

	"github.com/jcelliott/lumber"

	"github.com/nanopack/mist/clients"
	"github.com/nanopack/mist/server"
)

var (
	testAddr = "127.0.0.1:2445"
	testTag  = "hello"
	testMsg  = "world"
)

// TestMain
func TestMain(m *testing.M) {
	lumber.Level(lumber.LvlInt("fatal"))

	server.StartTCP(testAddr, nil)

	//
	os.Exit(m.Run())
}

// TestTCPClientConnect tests to ensure a client can connect to a running server
func TestTCPClientConnect(t *testing.T) {
	client, err := clients.New(testAddr, "")
	if err != nil {
		t.Fatalf("Client failed to connect - %v", err.Error())
		t.FailNow()
	}
	defer client.Close()

	//
	if err := client.Ping(); err != nil {
		t.Fatalf("ping failed")
	}
	if msg := <-client.Messages(); msg.Data != "pong" {
		t.Fatalf("Unexpected data: Expecting 'pong' got %s", msg.Data)
	}
	client.Ping()
}

// TestBadTCPClientConnect tests to ensure a client can connect to a running server
func TestBadTCPClientConnect(t *testing.T) {
	client, err := clients.New("321321321", "")
	if err == nil {
		t.Fatalf("Client succeeded to connect")
	}
	client, err = clients.New(testAddr, "hil")
	if err != nil {
		t.Fatalf("Client failed to connect - %v", err.Error())
		t.FailNow()
	}
	defer client.Close()

	//
	if err := client.Ping(); err != nil {
		t.Fatalf("ping failed")
	}
	if msg := <-client.Messages(); msg.Data != "pong" {
		t.Fatalf("Unexpected data: Expecting 'pong' got %s", msg.Data)
	}
	client.Ping()
}

// TestTCPClient tests to ensure a client can run all of its expected commands;
// we don't have to actually test any of the results of the commands since those
// are already tested in other tests (proxy_test and subscriptions_test in the
// core package)
func TestTCPClient(t *testing.T) {

	//
	client, err := clients.New(testAddr, "")
	if err != nil {
		t.Fatalf("failed to connect - %v", err.Error())
		t.FailNow()
	}
	defer client.Close()

	// subscribe should fail with no tags
	if err := client.Subscribe([]string{}); err == nil {
		t.Fatalf("Subscription succeeded with missing tags!")
	}

	// test ability to subscribe
	if err := client.Subscribe([]string{"a"}); err != nil {
		t.Fatalf("client subscriptions failed %v", err.Error())
	}

	// test ability to list (subscriptions)
	if err := client.List(); err != nil {
		t.Fatalf("listing subscriptions failed %v", err.Error())
	}
	if msg := <-client.Messages(); msg.Data == "\"a\"" {
		t.Fatalf("Failed to 'list' - '%v' '%#v'", msg.Error, msg.Data)
	}

	// test publish
	if err := client.Publish([]string{"a"}, "testpublish"); err != nil {
		t.Fatalf("publishing failed %v", err.Error())
	}
	if err := client.Publish([]string{}, "nopublish"); err == nil {
		t.Fatalf("publishing no tags succeeded %v", err.Error())
	}
	if err := client.Publish([]string{"a"}, ""); err == nil {
		t.Fatalf("publishing no data succeeded %v", err.Error())
	}

	// test PublishAfter
	if err := client.PublishAfter([]string{"a"}, "testpublish", time.Second); err != nil {
		t.Fatalf("publishing failed %v", err.Error())
	}
	time.Sleep(time.Millisecond * 1500)

	// test ability to unsubscribe
	if err := client.Unsubscribe([]string{"a"}); err != nil {
		t.Fatalf("client unsubscriptions failed %v", err.Error())
	}
	if err := client.Unsubscribe([]string{}); err == nil {
		t.Fatalf("client unsubscriptions no tags succeeded %v", err.Error())
	}

	// test ability to list (no subscriptions)
	if err := client.List(); err != nil {
		t.Fatalf("listing subscriptions failed %v", err.Error())
	}
	if msg := <-client.Messages(); msg.Data != "" {
		t.Fatalf("Failed to 'list' - %v", msg.Error)
	}
}
