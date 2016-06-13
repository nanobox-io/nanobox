package ipControl_test

import (
	"net"
	"os"
	"testing"

	"github.com/nanobox-io/nanobox/util/ipControl"
)

// TestMain ...
func TestMain(m *testing.M) {
	ipControl.Flush()
	os.Exit(m.Run())
}

// TestReservingIps ...
func TestReservingIps(t *testing.T) {
	ipOne, err := ipControl.ReserveGlobal()
	if err != nil {
		t.Errorf("unable to reserve ip", err)
	}
	ipTwo, err := ipControl.ReserveGlobal()
	if err != nil {
		t.Errorf("unable to reserve ip", err)
	}
	ipThree, err := ipControl.ReserveLocal()
	if err != nil {
		t.Errorf("unable to reserve ip", err)
	}
	if ipOne.String() != "192.168.99.50" || ipTwo.String() != "192.168.99.51" || ipThree.String() != "192.168.0.50" {
		t.Errorf("incorrect ip addresses", ipOne, ipTwo, ipThree)
	}
}

// TestReturnIP ...
func TestReturnIP(t *testing.T) {
	err := ipControl.ReturnIP(net.ParseIP("192.168.99.50"))
	if err != nil {
		t.Errorf("unable to return ip", err)
	}
	err = ipControl.ReturnIP(net.ParseIP("192.168.99.51"))
	if err != nil {
		t.Errorf("unable to return ip", err)
	}
	err = ipControl.ReturnIP(net.ParseIP("192.168.0.50"))
	if err != nil {
		t.Errorf("unable to return ip", err)
	}
}

// TestReuseIP ...
func TestReuseIP(t *testing.T) {
	one, err := ipControl.ReserveGlobal()
	if err != nil {
		t.Errorf("unable to reserve ip", err)
	}
	ipTwo, err := ipControl.ReserveGlobal()
	if err != nil {
		t.Errorf("unable to reserve ip", err)
	}
	three, err := ipControl.ReserveLocal()
	if err != nil {
		t.Errorf("unable to reserve ip", err)
	}
	err = ipControl.ReturnIP(ipTwo)
	if err != nil {
		t.Errorf("unable to return ip", err)
	}
	ipTwoAgain, err := ipControl.ReserveGlobal()
	if err != nil {
		t.Errorf("unable to reserve ip", err)
	}
	if !ipTwo.Equal(ipTwoAgain) {
		t.Errorf("i should ahve recieved a repeat of %s but i got %s", ipTwo.String(), ipTwoAgain.String())
	}
	ipControl.ReturnIP(one)
	ipControl.ReturnIP(three)
}
