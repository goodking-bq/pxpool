package scanner

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"pxpool/models"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// Scanner 扫描器
type Scanner struct {
	IPcount   models.Int64Counter
	ScanCount models.Int64Counter
	Doing     models.Int64Counter
	DoChan    chan bool
	Chan      chan Address
	Wg        sync.WaitGroup
	logger    *logrus.Logger
}

// NewScanner xxxx
func NewScanner(mc int64, logger *logrus.Logger) *Scanner {
	return &Scanner{
		Chan:   make(chan Address, mc),
		Wg:     sync.WaitGroup{},
		logger: logger,
		DoChan: make(chan bool, mc),
	}
}

// ScanIP 扫描 ip端口
func (scanner *Scanner) ScanIP(ipc Address, dataChan chan *models.Proxy) {
	defer scanner.Done(ipc)
	// time.AfterFunc(2*time.Second, func() {
	// 	scanner.Done()
	// 	fmt.Printf("scan %s:%d canceled", ipc.IP, ipc.Port)
	// 	return
	// })
	var nip net.IP
	err := nip.UnmarshalText([]byte(ipc.IP))
	if err != nil {
		log.Println(err)
		return
	}
	d := &net.Dialer{Timeout: 500 * time.Millisecond}
	target := fmt.Sprintf("%s:%d", ipc.IP, ipc.Port)
	scanner.logger.Debugf("开始连接：%s:%d", ipc.IP, ipc.Port)
	conn, err := d.Dial("tcp", target)
	if err != nil {
		return
	}
	conn.Close()
	scanner.logger.Debugf("连接成功：%s:%d， 正在检查 ..", ipc.IP, ipc.Port)
	px := &models.Proxy{Host: ipc.IP, Port: fmt.Sprintf("%d", ipc.Port), Category: "http"}
	if models.CheckProxy(px) {
		scanner.logger.Infof("发现新的代理：%s", px.URL())
		dataChan <- px
		return
	}
	if models.CheckProxy(&models.Proxy{Host: ipc.IP, Port: fmt.Sprintf("%d", ipc.Port), Category: "socks5"}) {
		scanner.logger.Infof("发现新的代理：%s", px.URL())
		dataChan <- px
		return
	}
	scanner.logger.Debugf("%s:%d open ,and is not  proxy \n", ipc.IP, ipc.Port)
}

// ScanCidr 扫描ip段
func (scanner *Scanner) FromCidr(cidr []byte, dataChan chan *models.Proxy) {
	ads := NewAddresses()
	err := ads.UnmarshalCidrText(cidr)
	if err != nil {
		fmt.Println(err)
	}
	scanner.IPcount.Inc(int64(len(*ads)))
	for _, address := range *ads {
		// ScanChan <- address
		go scanner.ScanIP(address, dataChan)
	}
}

// FromFile 扫描ip文件
func (scanner *Scanner) FromFile(file io.Reader, dataChan chan *models.Proxy) error {
	fscanner := bufio.NewScanner(file)
	for fscanner.Scan() {
		cidr := fscanner.Text()
		log.Printf("扫描段 ： %s", cidr)
		ip, ipnet, err := net.ParseCIDR(cidr)
		if err != nil {
			scanner.logger.Errorf("文件行错误: %s", cidr)
		}
		for IP := ip.Mask(ipnet.Mask); ipnet.Contains(IP); inc(IP) {
			if IP.String() == "172.0.0.254" {
				log.Println(IP)
			}

			if IP.String() != ipnet.IP.String() && IP.String() != ipnet.Mask.String() {
				for _, port := range PORTS {
					addr := Address{IP: IP.String(), Port: port}
					scanner.Chan <- addr
					scanner.IPcount.Inc(1)
				}
			}

		}
	}

	if err := fscanner.Err(); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// Scan 开始扫描
func (scanner *Scanner) Scan(ctx context.Context, maxc int64, dataChan chan *models.Proxy) error {
	ticker := time.NewTicker(time.Second * 10)
	var lastTick int64

	for {
		for scanner.Doing.Get() < maxc {
			select {
			case <-ticker.C:
				scanner.logger.Infof("ip数：%d,正在执行: %d,已完成: %d\n", scanner.IPcount.Get(), scanner.Doing.Get(), scanner.ScanCount.Get())
				if scanner.ScanCount.Get() != lastTick {
					lastTick = scanner.ScanCount.Get()
				} else {
					scanner.logger.Fatalln("扫描意外停止停止。")
				}
			case address := <-scanner.Chan:
				scanner.Wg.Add(1)
				go scanner.ScanIP(address, models.ProxyChan)
				scanner.Doing.Inc(1)
			case <-ctx.Done():
				close(models.ProxyChan)
				break
			}
		}
		scanner.Wg.Wait()
	}
}

func (scanner *Scanner) Done(addr Address) {
	scanner.ScanCount.Inc(1)
	scanner.DoChan <- true
	scanner.Doing.Dec(1)
	scanner.Wg.Done()
	scanner.logger.WithFields(logrus.Fields{
		"ip":   addr.IP,
		"port": addr.Port,
	}).Debugln("扫描完成")
}
