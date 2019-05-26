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
	Doing     int
	DoChan    chan bool
	Chan      chan Address
}

func NewScanner(config *models.Config) *Scanner {
	return &Scanner{
		Chan:   make(chan Address, config.Scanner.MaxConcurrency),
		DoChan: make(chan bool, config.Scanner.MaxConcurrency),
	}
}

// ScanIP 扫描 ip端口
func (scanner *Scanner) ScanIP(ipc Address, dataChan chan *models.Proxy) {
	log.Printf("scan %s:%d", ipc.IP, ipc.Port)
	var nip net.IP
	err := nip.UnmarshalText([]byte(ipc.IP))
	if err == nil {
		now := time.Now()
		after := now.Add(time.Duration(1) * time.Second)
		target := fmt.Sprintf("%s:%d", ipc.IP, ipc.Port)
		d := &net.Dialer{Timeout: 500 * time.Millisecond, Deadline: after}
		log.Printf("scan %s:%d 11111", ipc.IP, ipc.Port)
		conn, err := d.Dial("tcp", target)
		log.Printf("scan %s:%d 22222", ipc.IP, ipc.Port)
		if err == nil {
			conn.Close()
			px := &models.Proxy{Host: ipc.IP, Port: fmt.Sprintf("%d", ipc.Port), Category: "http"}
			if models.CheckProxy(px) {
				fmt.Printf("%s:%d open ,and is a http proxy \n", ipc.IP, ipc.Port)
				dataChan <- px
				return
			}
			if models.CheckProxy(&models.Proxy{Host: ipc.IP, Port: fmt.Sprintf("%d", ipc.Port), Category: "socks5"}) {
				fmt.Printf("%s:%d open ,and is a socks5 proxy \n", ipc.IP, ipc.Port)
				dataChan <- px
				return
			}
			fmt.Printf("%s:%d open ,and is not  proxy \n", ipc.IP, ipc.Port)
		}
	}
	time.Sleep(500 * time.Millisecond)
	scanner.ScanCount++
	<-scanner.DoChan
	scanner.Doing--
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
		for IP := ip.Mask(ipnet.Mask); IP.String() != ip.String() && ipnet.Contains(IP) && IP.String() != ipnet.Mask.String(); inc(IP) {
			for _, port := range PORTS {
				addr := Address{IP: IP.String(), Port: port}
				scanner.Chan <- addr
				scanner.IPcount++
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
	for {
		for scanner.Doing < config.Scanner.MaxConcurrency {
			select {
			// case <-scanner.DoChan:
			// 	address := <-scanner.Chan
			// 	scanner.WaitGroup.Add(1)
			// 	go scanner.ScanIP(address, dataChan)
			// 	scanner.Doing++
			// 	fmt.Printf("总IP: %d,已完成: %d\n", scanner.IPcount, scanner.ScanCount)
			// 	break
			case address := <-scanner.Chan:
				go scanner.ScanIP(address, dataChan)
				scanner.Doing++
				scanner.DoChan <- true
				//fmt.Println("\033[H\033[2J")
				fmt.Printf("发现IP: %d,已完成: %d\n", scanner.IPcount, scanner.ScanCount)
				break
			case <-ctx.Done():
				close(dataChan)
				break
			}
		}
		//time.Sleep(10 * time.Millisecond)
		//fmt.Println(scanner.Doing, scanner.IPcount, scanner.ScanCount)
		//scanner.WaitGroup.Wait()

	}
}
