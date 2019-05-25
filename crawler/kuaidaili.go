package crawler

import (
	"errors"
	"log"
	"net/http"
	"pxpool/models"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// KdlSpider 快代理
type KdlSpider struct {
	Name     string
	DataChan chan *models.Proxy
}

// NewKdlSpider 创建
func NewKdlSpider(DataChan chan *models.Proxy) *KdlSpider {
	return &KdlSpider{Name: "kdl", DataChan: DataChan}
}

// ToSpider 创建
func (c *KdlSpider) ToSpider() *Spider {
	var i Spider
	i = c
	return &i
}

// Start 快代理爬虫
func (c *KdlSpider) Start() {
	log.Println("快代理爬虫 开始运行 ...")
	for _, url := range c.GetUrls() {
		err := c.Run(url)
		if err != nil {
			log.Println(err)
		}
		time.Sleep(3 * time.Second)
	}
	log.Println("快代理爬虫 运行结束")
}

// GetUrls 链接    GetUrls() []string
func (c *KdlSpider) GetUrls() []string {
	var urls []string
	for i := 1; i < 3; i++ {
		url := "https://www.kuaidaili.com/free/inha/" + strconv.Itoa(i) + "/"
		urls = append(urls, url)
	}
	return urls
}

// Run 抓起页面
func (c *KdlSpider) Run(url string) error {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, strings.NewReader(""))
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows; U; Windows NT 6.1; en-US) AppleWebKit/532.5 (KHTML, like Gecko) Chrome/4.0.249.0 Safari/532.5")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return errors.New("status code error: " + resp.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return err
	}
	doc.Find("tbody tr").Each(func(i int, s *goquery.Selection) {
		proxy := models.NewProxy(s.Find("td[data-title='IP']").First().Text(),
			s.Find("td[data-title='PORT']").First().Text(),
			strings.ToLower(s.Find("td[data-title='类型']").First().Text()),
		)
		c.DataChan <- proxy
	})
	return nil
}

func (c *KdlSpider) GetName() string {
	return c.Name
}

func (c *KdlSpider) SetName(n string) error {
	c.Name = n
	return nil
}
