package crawl

import (
	"time"
)

type crawl struct {
	Name string
}

// Crawl 爬虫接口
type Crawl interface {
	Run(url string) error
	Start()
	GetUrls() []string
	GetName() string
	SetName(n string) error
	ToCrawl() *Crawl
}

// Manager 爬虫管理器
type Manager struct {
	crawls map[string]*Crawl
}

// NewManager 创建Manager
func NewManager() *Manager {
	return &Manager{crawls: make(map[string]*Crawl)}
}

// NewDefaultManager 创建Manager
func NewDefaultManager() *Manager {
	crawls := make(map[string]*Crawl)
	kdl := NewKdlCrawl().ToCrawl()
	crawls[(*kdl).GetName()] = kdl
	return &Manager{crawls: crawls}
}

// Add 添加一个新爬虫
func (cm *Manager) Add(c *Crawl) error {
	cm.crawls[(*c).GetName()] = c
	return nil
}

// Start 开始所有爬虫
func (cm *Manager) Start() {
	for _, crawl := range cm.crawls {
		go (*crawl).Start()
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
				cm.Start()
			case stop := <-stopChan:
				if stop {
					return
				}
			}
		}
	}(crawlTicker)
	return stopChan
}

// StartAndTicker 开始并定时执行
func (cm *Manager) StartAndTicker() chan bool {
	cm.Start()
	return cm.StartTicker()
}
