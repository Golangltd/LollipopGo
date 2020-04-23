package net

import (
	"net"
)

// IPv4List 获取本机的 ipv4 列表.
func IPv4List() ([]net.IP, error) {
	itfs, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	var (
		itf      net.Interface
		addrs    []net.Addr
		addr     net.Addr
		ipNet    *net.IPNet
		ok       bool
		ipv4     net.IP
		ipv4List []net.IP
	)
	for _, itf = range itfs {
		if itf.Flags&net.FlagUp == 0 {
			continue
		}
		addrs, err = itf.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr = range addrs {
			ipNet, ok = addr.(*net.IPNet)
			if !ok || ipNet.IP.IsLoopback() {
				continue
			}
			ipv4 = ipNet.IP.To4()
			if ipv4 == nil {
				continue
			}
			ipv4List = append(ipv4List, ipv4)
		}
	}
	return ipv4List, nil
}
