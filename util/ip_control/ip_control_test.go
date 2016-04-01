package ip_control_test

import (
	"testing"
	"os"
	"github.com/nanobox-io/nanobox/util/ip_control"
	"net"
)


func TestMain(m *testing.M) {
	ip_control.Flush()
	os.Exit(m.Run())
}

func TestReservingIps(t *testing.T) {
	ipOne, err := ip_control.ReserveGlobal()
	if err != nil {
		t.Errorf("unable to reserve ip", err)
	}
	ipTwo, err := ip_control.ReserveGlobal()
	if err != nil {
		t.Errorf("unable to reserve ip", err)
	}
	ipThree, err := ip_control.ReserveLocal()
	if err != nil {
		t.Errorf("unable to reserve ip", err)
	}
	if ipOne.String() != "192.168.99.50" || ipTwo.String() != "192.168.99.51" || ipThree.String() != "192.168.0.50" {
		t.Errorf("incorrect ip addresses", ipOne, ipTwo, ipThree)
	}
}

func TestReturnIP(t *testing.T) {
	err := ip_control.ReturnIP(net.ParseIP("192.168.99.50"))
	if err != nil {
		t.Errorf("unable to return ip", err)
	}
	err = ip_control.ReturnIP(net.ParseIP("192.168.99.51"))
	if err != nil {
		t.Errorf("unable to return ip", err)
	}
	err = ip_control.ReturnIP(net.ParseIP("192.168.0.50"))
	if err != nil {
		t.Errorf("unable to return ip", err)
	}
}

func TestReuseIP(t *testing.T) {
	one, err := ip_control.ReserveGlobal()
	if err != nil {
		t.Errorf("unable to reserve ip", err)
	}
	ipTwo, err := ip_control.ReserveGlobal()
	if err != nil {
		t.Errorf("unable to reserve ip", err)
	}
	three, err := ip_control.ReserveLocal()
	if err != nil {
		t.Errorf("unable to reserve ip", err)
	}
	err = ip_control.ReturnIP(ipTwo)
	if err != nil {
		t.Errorf("unable to return ip", err)
	}
	ipTwoAgain, err := ip_control.ReserveGlobal()
	if err != nil {
		t.Errorf("unable to reserve ip", err)
	}
	if !ipTwo.Equal(ipTwoAgain) {
		t.Errorf("i should ahve recieved a repeat of %s but i got %s", ipTwo.String(), ipTwoAgain.String())
	}
	ip_control.ReturnIP(one)
	ip_control.ReturnIP(three)

}

