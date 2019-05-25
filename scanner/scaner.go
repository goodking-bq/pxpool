package scanner

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"pxpool/models"
	"time"
)

// Scanner 扫描器
type Scanner struct {
	IPcount   int
	ScanCount int
	Chan      chan Address
}

func NewScanner() *Scanner {
	return &Scanner{
		Chan: make(chan Address, 2),
	}
}

// ScanIP 扫描 ip端口
func (scanner *Scanner) ScanIP(ipc Address, dataChan chan *models.Proxy) {
	var nip net.IP
	nip.UnmarshalText([]byte(ipc.IP))
	d := net.Dialer{Timeout: 5 * time.Second}
	for _, port := range ipc.Ports {
		tcpaddr := &net.TCPAddr{IP: nip, Port: port}
		_, err := d.Dial("tcp", tcpaddr.String())
		if err != nil {
			fmt.Printf("%s:%d close\n", ipc.IP, port)
		} else {
			px := &models.Proxy{Host: ipc.IP, Port: fmt.Sprintf("%d", port), Category: "http"}
			if models.CheckProxy(px) {
				fmt.Printf("%s:%d open ,and is a http proxy \n", ipc.IP, port)
				dataChan <- px
				continue
			}
			if models.CheckProxy(&models.Proxy{Host: ipc.IP, Port: fmt.Sprintf("%d", port), Category: "socks5"}) {
				fmt.Printf("%s:%d open ,and is a socks5 proxy \n", ipc.IP, port)
				dataChan <- px
				continue
			}
			fmt.Printf("%s:%d open ,and is not  proxy \n", ipc.IP, port)
		}
	}
	scanner.ScanCount++
}

// ScanCidr 扫描ip段
func (scanner *Scanner) ScanCidr(cidr []byte, dataChan chan *models.Proxy) {
	ads := NewAddresses()
	err := ads.UnmarshalCidrText(cidr)
	if err != nil {
		fmt.Println(err)
	}
	scanner.IPcount += len(*ads)
	for _, address := range *ads {
		// ScanChan <- address
		go scanner.ScanIP(address, dataChan)
	}
}

// ScanFile 扫描ip文件
func (scanner *Scanner) ScanFile(f string, dataChan chan *models.Proxy) error {
	file, err := os.Open(f)
	if err != nil {
		log.Println(err)
		return err
	}
	defer file.Close()

	fscanner := bufio.NewScanner(file)
	for fscanner.Scan() {
		ip, ipnet, err := net.ParseCIDR(fscanner.Text())
		if err != nil {
			return err
		}
		i := 1
		for IP := ip.Mask(ipnet.Mask); IP.String() != ip.String() && ipnet.Contains(IP) && IP.String() != ipnet.Mask.String(); inc(IP) {
			addr := Address{IP: IP.String(), Ports: PORTS}
			scanner.Chan <- addr
			scanner.IPcount++
			func() {
				address := <-scanner.Chan
				(*scanner).ScanIP(address, dataChan)
			}()
			i++
			if i > 100 {
				break
			}
		}
		return nil

	}

	if err := fscanner.Err(); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// Scan 开始扫描
func (scanner *Scanner) Scan(ctx context.Context, config *models.Config, dataChan chan *models.Proxy) error {
	if config.Scanner.File != "" {
		go scanner.ScanFile(config.Scanner.File, dataChan)
	} else if config.Scanner.Cidr != "" {
		scanner.ScanCidr([]byte(config.Scanner.Cidr), dataChan)
	} else {
		return errors.New("未给扫描目标")
	}
	return nil
}
