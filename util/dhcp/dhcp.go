package dhcp

import (
	"errors"
	"net"
	"sync"

	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/data"
	"github.com/nanobox-io/nanobox/util/locker"
)

var (
	errIPNotFound = errors.New("Ip Not Found")
	mutex         = sync.Mutex{}
)

// IPSpace ...
type IPSpace struct {
	GlobalIP  net.IP
	GlobalNet net.IPNet
	LocalIP   net.IP
	LocalNet  net.IPNet
}

// ReserveGlobal ...
func ReserveGlobal() (net.IP, error) {

	locker.GlobalLock()
	defer locker.GlobalUnlock()
	mutex.Lock()
	defer mutex.Unlock()

	//
	ipSpace, err := getIPSpace()
	if err != nil {
		return nil, err
	}

	//
	reservedIPs, err := getReserved()
	if err != nil {
		return nil, err
	}

	//
	for ip := ipSpace.GlobalIP; ipSpace.GlobalNet.Contains(ip); inc(ip) {
		if !contains(reservedIPs, ip) {
			setReserved(append(reservedIPs, ip))
			if err != nil {
				return nil, err
			}
			return ip, nil
		}
	}

	return nil, errIPNotFound
}

// Flush ...
func Flush() {

	locker.GlobalLock()
	defer locker.GlobalUnlock()
	mutex.Lock()
	defer mutex.Unlock()

	data.Delete("global", "ipreserved")
	data.Delete("global", "ipreserved")
}

// ReserveLocal ...
func ReserveLocal() (net.IP, error) {

	locker.GlobalLock()
	defer locker.GlobalUnlock()
	mutex.Lock()
	defer mutex.Unlock()

	//
	ipSpace, err := getIPSpace()
	if err != nil {
		return nil, err
	}

	//
	reservedIPs, err := getReserved()
	if err != nil {
		return nil, err
	}

	//
	for ip := ipSpace.LocalIP; ipSpace.LocalNet.Contains(ip); inc(ip) {
		if !contains(reservedIPs, ip) {
			setReserved(append(reservedIPs, ip))
			if err != nil {
				return nil, err
			}
			return ip, nil
		}
	}

	return nil, errIPNotFound
}

// ReturnIP ...
func ReturnIP(ip net.IP) error {

	locker.GlobalLock()
	defer locker.GlobalUnlock()
	mutex.Lock()
	defer mutex.Unlock()

	//
	reservedIPs, err := getReserved()
	if err != nil {
		return err
	}

	//
	for i, reservedIP := range reservedIPs {
		if reservedIP.Equal(ip) {
			return setReserved(append(reservedIPs[:i], reservedIPs[i+1:]...))
		}
	}

	return nil
}

// getIPSpace do not store the space on the disk.
func getIPSpace() (IPSpace, error) {
	ipSpace := IPSpace{}

	// there was no data stored for ip space so we need to populate it
	ip, ipNet, err := net.ParseCIDR(config.Viper().GetString("external-network-space"))
	if err != nil {
		return ipSpace, err
	}
	ipSpace.GlobalIP = ip
	ipSpace.GlobalNet = *ipNet

	//
	ip, ipNet, err = net.ParseCIDR(config.Viper().GetString("internal-network-space"))
	if err != nil {
		return ipSpace, err
	}
	ipSpace.LocalIP = ip
	ipSpace.LocalNet = *ipNet

	return ipSpace, nil
}

// contains ...
func contains(ips []net.IP, ip net.IP) bool {

	//
	for _, setIP := range ips {
		if setIP.Equal(ip) {
			return true
		}
	}
	return false
}

// getReserved ...
func getReserved() ([]net.IP, error) {
	ips := []net.IP{}
	data.Get("global", "ipreserved", &ips)
	return ips, nil
}

// setReserved ...
func setReserved(ips []net.IP) error {
	return data.Put("global", "ipreserved", ips)
}

// inc ...
func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}
