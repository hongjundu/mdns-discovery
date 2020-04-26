package utils

import (
	"net"
)

func IpAddress() (ip net.IP, err error) {

	interfaceAddrs, err := net.InterfaceAddrs()

	if err != nil {
		return
	}

	for _, address := range interfaceAddrs {
		// check the address type and if it is not a loopback the display it
		// = GET LOCAL IP ADDRESS
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ip = ipnet.IP
				break
			}
		}
	}

	return
}
