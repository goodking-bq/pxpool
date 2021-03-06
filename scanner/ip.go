package scanner

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

// PORTS 默认端口
var PORTS = []int{80, 8080, 3128, 8081, 9080, 10808}

// Address 扫描的ip
type Address struct {
	IP   string
	Port int
}

// Addresses 列表
type Addresses []Address

// NewAddresses 新建
func NewAddresses() *Addresses {
	return &Addresses{}
}

// SetPorts 设置端口
func (ads *Addresses) SetPorts(port int) {
	for _, address := range *ads {
		address.Port = port
	}
}

// UnmarshalCidrText 返回 cidr所有的ip
func (ads *Addresses) UnmarshalCidrText(cidr []byte) error {
	ip, ipnet, err := net.ParseCIDR(string(cidr))
	if err != nil {
		return err
	}
	for IP := ip.Mask(ipnet.Mask); ipnet.Contains(IP); inc(IP) {
		//IP := net.IP([]byte(ip.String()))
		for port := range PORTS {
			addr := Address{IP: IP.String(), Port: port}
			*ads = append(*ads, addr)
		}

	}
	*ads = (*ads)[len(PORTS) : len(*ads)-len(PORTS)]
	return nil
}

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

// UnmarshalCidrFile 从文件读取
func (ads *Addresses) UnmarshalCidrFile(f string) error {
	file, err := os.Open(f)
	if err != nil {
		log.Println(err)
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
		ads.UnmarshalCidrText(scanner.Bytes())
	}

	if err := scanner.Err(); err != nil {
		log.Println(err)
		return err
	}
	return nil
}
