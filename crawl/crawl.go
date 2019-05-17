package crawl

import (
	"errors"
	"math/rand"
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
	return stopChan
}
