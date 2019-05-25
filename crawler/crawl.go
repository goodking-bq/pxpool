package crawler

import (
	"log"
	"pxpool/models"
	"pxpool/storage"
	"time"
)

// Spider 爬虫接口
type Spider interface {
	Run(url string) error
	Start() // 开始怕去县城
	GetUrls() []string
	GetName() string
	SetName(n string) error
	ToSpider() *Spider
}

// Crawl 爬虫管理器
type Crawl struct {
	DataChan   chan *models.Proxy // 数据交换用
	ExitSignal chan bool          // 退出信号
	Spiders    map[string]*Spider
	storage    *storage.Storager
}

// NewCrawl 创建Crawl
func NewCrawl(storage *storage.Storager) *Crawl {
	return &Crawl{
		Spiders:    make(map[string]*Spider),
		storage:    storage,
		DataChan:   make(chan *models.Proxy),
		ExitSignal: make(chan bool),
	}
}

// NewDefaultCrawl 创建Crawl
func NewDefaultCrawl(storage *storage.Storager, DataChan chan *models.Proxy) *Crawl {
	spiders := make(map[string]*Spider)
	kdl := NewKdlSpider(DataChan).ToSpider()
	spiders[(*kdl).GetName()] = kdl
	return &Crawl{
		Spiders:    spiders,
		storage:    storage,
		DataChan:   DataChan,
		ExitSignal: make(chan bool),
	}
}

// Add 添加一个新爬虫
func (cm *Crawl) Add(c *Spider) error {
	cm.Spiders[(*c).GetName()] = c
	return nil
}

// Save 保存
func (cm *Crawl) Save() { // 开始接受
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
func (cm *Crawl) Crawl() {
	for _, spider := range cm.Spiders {
		go (*spider).Start()
	}
}

func (cm *Crawl) Start() {
	cm.Crawl()
	cm.Save()
}

// StartTicker 开始爬虫循环跑
func (cm *Crawl) StartTicker() {
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
func (cm *Crawl) StartAndTicker() {
	cm.Crawl()
	cm.StartTicker()
	cm.Save()
}
