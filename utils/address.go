package utils

import "net"

func GetLocalIP() []string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil
	}
	ret := make([]string, 0)
	for _, address := range addrs {
		if ipNet, ok := address.(*net.IPNet); ok {
			if ipNet.IP.To4() == nil {
				continue
			}
			if ipNet.IP.IsGlobalUnicast() {
				ret = append(ret, ipNet.IP.String())
			}

		}
	}
	return ret
}
