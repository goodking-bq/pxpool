package crawler

import (
	"pxpool/models"
	"time"

	"github.com/sirupsen/logrus"
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
	logger     *logrus.Logger
	config     *models.Config
}

// NewCrawl 创建Crawl
func NewCrawl(logger *logrus.Logger, config *models.Config, DataChan chan *models.Proxy) *Crawl {
	return &Crawl{
		Spiders:    make(map[string]*Spider),
		DataChan:   DataChan,
		ExitSignal: make(chan bool),
		logger:     logger,
		config:     config,
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
	if cm.config.Crawl.Ticker > 0 {
		cm.StartTicker()
	}
}

// StartTicker 开始爬虫循环跑
func (cm *Crawl) StartTicker() {
	cm.logger.Debugln("爬虫60秒后再次运行")
	crawlTicker := time.NewTicker(time.Second * time.Duration(cm.config.Crawl.Ticker))
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
