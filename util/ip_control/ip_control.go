package ip_control

import (
	"errors"
	"github.com/nanobox-io/nanobox/util/data"
	"github.com/nanobox-io/nanobox/util/locker"
	"github.com/nanobox-io/nanobox/util/config"
	"net"
	"sync"
)

type IPSpace struct {
	GlobalIP  net.IP
	GlobalNet net.IPNet
	LocalIP   net.IP
	LocalNet  net.IPNet
}

var IpNotFound = errors.New("Ip Not Found")
var mutex = sync.Mutex{}

func ReserveGlobal() (net.IP, error) {
	locker.GlobalLock()
	defer locker.GlobalUnlock()
	mutex.Lock()
	defer mutex.Unlock()

	ipSpace, err := getIpSpace()
	if err != nil {
		return nil, err
	}
	reservedIPs, err := getReserved()
	if err != nil {
		return nil, err
	}

	for ip := ipSpace.GlobalIP; ipSpace.GlobalNet.Contains(ip); inc(ip) {
		if !contains(reservedIPs, ip) {
			setReserved(append(reservedIPs, ip))
			if err != nil {
				return nil, err
			}
			return ip, nil
		}
	}
	return nil, IpNotFound
}

func Flush() {
	locker.GlobalLock()
	defer locker.GlobalUnlock()
	mutex.Lock()
	defer mutex.Unlock()

	data.Delete("global", "ipreserved")
	data.Delete("global", "ipreserved")
}

func ReserveLocal() (net.IP, error) {
	locker.GlobalLock()
	defer locker.GlobalUnlock()
	mutex.Lock()
	defer mutex.Unlock()

	ipSpace, err := getIpSpace()
	if err != nil {
		return nil, err
	}
	reservedIPs, err := getReserved()
	if err != nil {
		return nil, err
	}
	for ip := ipSpace.LocalIP; ipSpace.LocalNet.Contains(ip); inc(ip) {
		if !contains(reservedIPs, ip) {
			setReserved(append(reservedIPs, ip))
			if err != nil {
				return nil, err
			}
			return ip, nil
		}
	}
	return nil, IpNotFound
}

func ReturnIP(ip net.IP) error {
	locker.GlobalLock()
	defer locker.GlobalUnlock()
	mutex.Lock()
	defer mutex.Unlock()

	reservedIPs, err := getReserved()
	if err != nil {
		return err
	}

	for i, reservedIP := range reservedIPs {
		if reservedIP.Equal(ip) {
			return setReserved(append(reservedIPs[:i], reservedIPs[i+1:]...))
		}
	}

	return nil
}

// do not store the space on the disk.
func getIpSpace() (IPSpace, error) {
	ipSpace := IPSpace{}

	// there was no data stored for ip space
	// so we need to populate it
	ip, ipNet, err := net.ParseCIDR(config.Viper().GetString("external-network-space"))
	if err != nil {
		return ipSpace, err
	}
	ipSpace.GlobalIP = ip
	ipSpace.GlobalNet = *ipNet

	ip, ipNet, err = net.ParseCIDR(config.Viper().GetString("internal-network-space"))
	if err != nil {
		return ipSpace, err
	}
	ipSpace.LocalIP = ip
	ipSpace.LocalNet = *ipNet

	return ipSpace, nil
}

func contains(ips []net.IP, ip net.IP) bool {
	for _, setIp := range ips {
		if setIp.Equal(ip) {
			return true
		}
	}
	return false
}

func getReserved() ([]net.IP, error) {
	ips := []net.IP{}
	data.Get("global", "ipreserved", &ips)
	return ips, nil
}

func setReserved(ips []net.IP) error {
	return data.Put("global", "ipreserved", ips)
}

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}
