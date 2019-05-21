package scanner

import (
	"net"
)

// PORTS 默认端口
var PORTS = []int{80, 8080, 3128, 8081, 9080, 10808}

// IPForScaner 扫描的ip
type IPForScaner struct {
	IP    string
	Ports []int
}

// UnmarshalCidr 返回 cidr所有的ip
func UnmarshalCidr(cidr string) ([]string, error) {
	var ips []string
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		ips = append(ips, ip.String())
	}
	return ips[1 : len(ips)-1], nil
}

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}
