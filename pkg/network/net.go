package network

import "net"

func CheckIpCidr(srcip string, cidrlist []string) bool {
	srcipcheck := net.ParseIP(srcip)
	for _, v := range cidrlist {
		_, ipnet, err := net.ParseCIDR(v)
		if err != nil {
			continue
		}
		if ipnet.Contains(srcipcheck) {
			return true
		}
	}

	return false
}
