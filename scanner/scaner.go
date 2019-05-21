package scanner

import (
	"fmt"
	"net"
	"time"

	"../model"
)

// ScanIP 扫描 ip端口
func ScanIP(ipc IPForScaner) {
	var nip net.IP
	nip.UnmarshalText([]byte(ipc.IP))
	d := net.Dialer{Timeout: 5 * time.Second}
	for _, port := range ipc.Ports {
		tcpaddr := &net.TCPAddr{IP: nip, Port: port}
		_, err := d.Dial("tcp", tcpaddr.String())
		if err != nil {
			fmt.Printf("%s:%d close\n", ipc.IP, port)
		} else {
			px := &model.Proxy{Ip: ipc.IP, Port: fmt.Sprintf("%d", port), Category: "http"}
			if model.CheckProxy(px) {
				fmt.Printf("%s:%d open ,and is a http proxy \n", ipc.IP, port)
				continue
			}
			if model.CheckProxy(&model.Proxy{Ip: ipc.IP, Port: fmt.Sprintf("%d", port), Category: "socks5"}) {
				fmt.Printf("%s:%d open ,and is a socks5 proxy \n", ipc.IP, port)
				continue
			}
			fmt.Printf("%s:%d open ,and is not  proxy \n", ipc.IP, port)

			// _, err := conn.Write([]byte("CONNECT baidu.com HTTP/1.1"))
			// if err != nil {
			// 	fmt.Printf("%s:%d open ,but not proxy \n", ipc.IP, port)
			// }
			// result, err := ioutil.ReadAll(conn)
			// if err != nil {
			// 	fmt.Printf("%s:%d open ,but not proxy \n", ipc.IP, port)
			// }
			// fmt.Println(string(result))
			// fmt.Printf("%s:%d open -----------------\n", ipc.IP, port)
		}
	}
}

// ScanCidr 扫描ip段
func ScanCidr(cidr string) {
	ips, _ := UnmarshalCidr(cidr)
	for _, ip := range ips {
		ipScan := IPForScaner{IP: ip, Ports: PORTS}
		go ScanIP(ipScan)
	}
}
