package crawl

import (
	"log"
	"time"

	"../model"
	"../storage"
)

type crawl struct {
	Name string
}

// Crawl 爬虫接口
type Crawl interface {
	Run(url string) error
	Start() // 开始怕去县城
	GetUrls() []string
	GetName() string
	SetName(n string) error
	ToCrawl() *Crawl
}

// Manager 爬虫管理器
type Manager struct {
	crawls     map[string]*Crawl
	storage    *storage.Storager
	DataChan   chan *model.Proxy // 数据交换用
	ExitSignal chan bool         // 退出信号
}

// NewManager 创建Manager
func NewManager(storage *storage.Storager) *Manager {
	return &Manager{crawls: make(map[string]*Crawl), storage: storage, DataChan: make(chan *model.Proxy), ExitSignal: make(chan bool)}
}

// NewDefaultManager 创建Manager
func NewDefaultManager(storage *storage.Storager) *Manager {
	crawls := make(map[string]*Crawl)
	DataChan := make(chan *model.Proxy)
	kdl := NewKdlCrawl(DataChan).ToCrawl()
	crawls[(*kdl).GetName()] = kdl
	return &Manager{crawls: crawls, storage: storage, DataChan: DataChan, ExitSignal: make(chan bool)}
}

// Add 添加一个新爬虫
func (cm *Manager) Add(c *Crawl) error {
	cm.crawls[(*c).GetName()] = c
	return nil
}

func (cm *Manager) Save() { // 开始接受
	for {
		select {
		case proxy := <-cm.DataChan:
			go (*cm.storage).AddOrUpdateProxy(proxy)
		case stop := <-cm.ExitSignal:
			if stop {
				close(cm.DataChan)
				return
			}
		}
	}

}

// Crawl 开始所有爬虫
func (cm *Manager) Crawl() {
	for _, crawl := range cm.crawls {
		go (*crawl).Start()
	}
}

func (cm *Manager) Start() {
	cm.Crawl()
	cm.Save()
}

// StartTicker 开始爬虫循环跑
func (cm *Manager) StartTicker() {
	log.Println("爬虫60秒后再次运行")
	crawlTicker := time.NewTicker(time.Second * 60)
	go func(ticker *time.Ticker) {
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				cm.Start()
			case stop := <-cm.ExitSignal:
				if stop {
					close(cm.DataChan)
					return
				}
			}
		}
	}(crawlTicker)
}

// StartAndTicker 开始并定时执行
func (cm *Manager) StartAndTicker() {
	cm.Crawl()
	cm.StartTicker()
	cm.Save()
}
