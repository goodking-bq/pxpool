package scanner

import (
	"fmt"
	"net"
)

// ScanIP 扫描 ip端口
func ScanIP(ip string, ports []int) {
	var nip net.IP
	nip.UnmarshalText([]byte(ip))
	for _, port := range ports {
		tcpaddr := &net.TCPAddr{IP: nip, Port: port}
		_, err := net.DialTCP("tcp", nil, tcpaddr)
		if err != nil {
			fmt.Printf("%s:%d close\n", ip, port)
		} else {
			fmt.Printf("%s:%d open\n", ip, port)
		}
	}
}
