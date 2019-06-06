package scanner

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
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
	config    *models.Config
}

// NewScanner xxxx
func NewScanner(config *models.Config, logger *logrus.Logger) *Scanner {
	return &Scanner{
		Chan:   make(chan Address, config.Scanner.MaxConcurrency),
		Wg:     sync.WaitGroup{},
		logger: logger,
		DoChan: make(chan bool, config.Scanner.MaxConcurrency),
		config: config,
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

// FromCidr 扫描ip段
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

// MakeAddress xxx
func (scanner *Scanner) MakeAddress() error {
	if scanner.config.Scanner.File != "" {
		file, err := os.Open(scanner.config.Scanner.File)
		if err != nil {
			scanner.logger.Errorf("文件 %s 不存在", scanner.config.Scanner.File)
			return err
		}
		defer file.Close()
		go scanner.FromFile(file, models.ProxyChan)
	} else if scanner.config.Scanner.Cidr != "" {
		go scanner.FromCidr([]byte(scanner.config.Scanner.Cidr), models.ProxyChan)
	} else {
		scanner.logger.Errorln("未给扫描目标")
		return errors.New("未给扫描目标")
	}
	return nil
}

// Scan scan
func (scanner *Scanner) Scan(ctx context.Context) error {
	ticker := time.NewTicker(time.Second * 60)
	var lastTick int64
	for i := int64(0); i < int64(scanner.config.Scanner.MaxConcurrency); i++ {
		scanner.DoChan <- true
	}
	for {
		//for scanner.Doing.Get() < maxConcurrency {
		select {
		case <-ticker.C:
			scanner.logger.Infof("ip数：%d,正在执行: %d,已完成: %d\n", scanner.IPcount.Get(), scanner.Doing.Get(), scanner.ScanCount.Get())
			if scanner.ScanCount.Get() != lastTick {
				lastTick = scanner.ScanCount.Get()
			} else {
				scanner.logger.Fatalln("扫描停止。")
			}
		case address := <-scanner.Chan:
			<-scanner.DoChan
			scanner.Wg.Add(1)
			go scanner.ScanIP(address, models.ProxyChan)
			scanner.Doing.Inc(1)

		case <-ctx.Done():
			close(models.ProxyChan)
			break
		}
		//}
	}
}

// Done done one
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
