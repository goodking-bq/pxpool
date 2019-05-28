package crawler

import (
	"log"
	"pxpool/models"
	"time"

	"github.com/spf13/cobra"
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
}

// NewCrawl 创建Crawl
func NewCrawl() *Crawl {
	return &Crawl{
		Spiders:    make(map[string]*Spider),
		DataChan:   make(chan *models.Proxy),
		ExitSignal: make(chan bool),
	}
}

// NewDefaultCrawl 创建Crawl
func NewDefaultCrawl(DataChan chan *models.Proxy) *Crawl {
	spiders := make(map[string]*Spider)
	kdl := NewKdlSpider(DataChan).ToSpider()
	spiders[(*kdl).GetName()] = kdl
	return &Crawl{
		Spiders:    spiders,
		DataChan:   DataChan,
		ExitSignal: make(chan bool),
	}
}

// Add 添加一个新爬虫
func (cm *Crawl) Add(c *Spider) error {
	cm.Spiders[(*c).GetName()] = c
	return nil
}

// Crawl 开始所有爬虫
func (cm *Crawl) Crawl() {
	for _, spider := range cm.Spiders {
		go (*spider).Start()
	}
}

func (cm *Crawl) Start() {
	cm.Crawl()
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
}

// Command 爬虫命令
func Command(cmd *cobra.Command) {}
