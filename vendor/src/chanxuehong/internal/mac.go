package internal

import (
	"bytes"
	"net"

	"github.com/chanxuehong/rand"
)

var MAC [6]byte = getMAC() // One MAC of this machine; Particular case, it is a random bytes.

var zeroMAC [8]byte

func getMAC() (mac [6]byte) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return genMAC()
	}

	// Gets a MAC from interfaces of this machine,
	// the MAC of up state interface is preferred.
	found := false // Says it has found a MAC
	for _, itf := range interfaces {
		if itf.Flags&net.FlagLoopback == net.FlagLoopback ||
			itf.Flags&net.FlagPointToPoint == net.FlagPointToPoint {
			continue
		}

		switch hardwareAddr := itf.HardwareAddr; len(hardwareAddr) {
		case 6: // MAC-48, EUI-48
			if bytes.Equal(hardwareAddr, zeroMAC[:6]) {
				continue
			}
			if itf.Flags&net.FlagUp == 0 {
				if !found {
					copy(mac[:], hardwareAddr)
					found = true
				}
				continue
			}
			copy(mac[:], hardwareAddr)
			return
		case 8: // EUI-64
			if bytes.Equal(hardwareAddr, zeroMAC[:]) {
				continue
			}
			if itf.Flags&net.FlagUp == 0 {
				if !found {
					copy(mac[:3], hardwareAddr)
					copy(mac[3:], hardwareAddr[5:])
					found = true
				}
				continue
			}
			copy(mac[:3], hardwareAddr)
			copy(mac[3:], hardwareAddr[5:])
			return
		}
	}
	if found {
		return
	}

	return genMAC()
}

// generates a random MAC.
func genMAC() (mac [6]byte) {
	rand.Read(mac[:])
	mac[0] |= 0x01 // multicast
	return
}
