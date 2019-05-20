package crawl

import (
	"errors"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

// Proxy 代理
type Proxy struct {
	Ip         string
	Port       string
	category   string
	joinTime   string
	verifyTime string
}

// URL 获取代理的地址
func (p *Proxy) URL() string {
	return p.category + "://" + p.Ip + ":" + p.Port
}

func (p *Proxy) check() bool {
	log.Printf("check %s", p.URL())
	proxy, _ := url.Parse(p.URL())
	netTransport := &http.Transport{
		//Proxy: http.ProxyFromEnvironment,
		Proxy: http.ProxyURL(proxy),
		Dial: func(netw, addr string) (net.Conn, error) {
			c, err := net.DialTimeout(netw, addr, time.Second*time.Duration(10))
			if err != nil {
				return nil, err
			}
			return c, nil
		},
		MaxIdleConnsPerHost:   10,                             //每个host最大空闲连接
		ResponseHeaderTimeout: time.Second * time.Duration(5), //数据收发5秒超时
	}
	client := &http.Client{
		Timeout:   time.Second * 10,
		Transport: netTransport}
	req, err := http.NewRequest("GET", "https://baidu.com", strings.NewReader(""))
	if err != nil {
		return true
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows; U; Windows NT 6.1; en-US) AppleWebKit/532.5 (KHTML, like Gecko) Chrome/4.0.249.0 Safari/532.5")
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	if resp.StatusCode != 200 {
		return false
	}
	return true
}

type proxysMap struct {
	sync.Map
}

// Proxys 所有的代理
var Proxys proxysMap

func (p *proxysMap) Random() (Proxy, error) {
	var ips []string
	Proxys.Range(func(k, _ interface{}) bool {
		ips = append(ips, k.(string))
		return true
	})
	l := len(ips)
	if l == 0 {
		return Proxy{}, errors.New("没有缓存代理")
	}
	n := rand.Intn(l)
	_p, _ := Proxys.Load(ips[n])
	return _p.(Proxy), nil
}
func (p *proxysMap) Check() {
	log.Println("开始检查 。。。")
	Proxys.Range(func(k, v interface{}) bool {
		px := v.(Proxy)
		isActive := px.check()
		if !isActive {
			Proxys.Delete(k)
		}
		return true
	})

}

// Crawl 爬虫接口
type Crawl interface {
	Run(url string) error
	Start()
	GetUrls() []string
}

// Manager 爬虫管理器
type Manager struct {
	crawls []Crawl
}

// Add 添加一个新爬虫
func (cm *Manager) Add(c *Crawl) error {
	cm.crawls = append(cm.crawls, *c)
	return nil
}

// Start 开始所有爬虫
func (cm *Manager) Start(ticker bool) {
	for _, crawl := range cm.crawls {
		go crawl.Start()
	}
	if ticker == true {
		cm.StartTicker()
	}
}

// StartTicker 开始爬虫循环跑
func (cm *Manager) StartTicker() chan bool {
	crawlTicker := time.NewTicker(time.Second * 60)
	checkTicker := time.NewTicker(time.Second * 60)

	stopChan := make(chan bool)
	go func(ticker *time.Ticker) {
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				cm.Start(false)
			case stop := <-stopChan:
				if stop {
					return
				}
			}
		}
	}(crawlTicker)

	go func(ticker *time.Ticker) {
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				Proxys.Check()
			case stop := <-stopChan:
				if stop {
					return
				}
			}
		}
	}(checkTicker)
	return stopChan
}
