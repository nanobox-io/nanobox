package dhcp_test

import (
	"net"
	"os"
	"testing"

	"github.com/nanobox-io/nanobox/util/dhcp"
)

// TestMain ...
func TestMain(m *testing.M) {
	dhcp.Flush()
	os.Exit(m.Run())
}

// TestReservingIps ...
func TestReservingIps(t *testing.T) {
	ipOne, err := dhcp.ReserveGlobal()
	if err != nil {
		t.Errorf("unable to reserve ip", err)
	}
	ipTwo, err := dhcp.ReserveGlobal()
	if err != nil {
		t.Errorf("unable to reserve ip", err)
	}
	ipThree, err := dhcp.ReserveLocal()
	if err != nil {
		t.Errorf("unable to reserve ip", err)
	}
	if ipOne.String() != "192.168.99.50" || ipTwo.String() != "192.168.99.51" || ipThree.String() != "192.168.0.50" {
		t.Errorf("incorrect ip addresses", ipOne, ipTwo, ipThree)
	}
}

// TestReturnIP ...
func TestReturnIP(t *testing.T) {
	err := dhcp.ReturnIP(net.ParseIP("192.168.99.50"))
	if err != nil {
		t.Errorf("unable to return ip", err)
	}
	err = dhcp.ReturnIP(net.ParseIP("192.168.99.51"))
	if err != nil {
		t.Errorf("unable to return ip", err)
	}
	err = dhcp.ReturnIP(net.ParseIP("192.168.0.50"))
	if err != nil {
		t.Errorf("unable to return ip", err)
	}
}

// TestReuseIP ...
func TestReuseIP(t *testing.T) {
	one, err := dhcp.ReserveGlobal()
	if err != nil {
		t.Errorf("unable to reserve ip", err)
	}
	ipTwo, err := dhcp.ReserveGlobal()
	if err != nil {
		t.Errorf("unable to reserve ip", err)
	}
	three, err := dhcp.ReserveLocal()
	if err != nil {
		t.Errorf("unable to reserve ip", err)
	}
	err = dhcp.ReturnIP(ipTwo)
	if err != nil {
		t.Errorf("unable to return ip", err)
	}
	ipTwoAgain, err := dhcp.ReserveGlobal()
	if err != nil {
		t.Errorf("unable to reserve ip", err)
	}
	if !ipTwo.Equal(ipTwoAgain) {
		t.Errorf("i should ahve recieved a repeat of %s but i got %s", ipTwo.String(), ipTwoAgain.String())
	}
	dhcp.ReturnIP(one)
	dhcp.ReturnIP(three)
}
