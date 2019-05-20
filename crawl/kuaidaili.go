package crawl

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// KdlCrawl 快代理
type KdlCrawl struct {
}

// Start 快代理爬虫
func (c *KdlCrawl) Start() {
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

// GetUrls 链接
func (c *KdlCrawl) GetUrls() []string {
	var urls []string
	for i := 1; i < 3; i++ {
		url := "https://www.kuaidaili.com/free/inha/" + strconv.Itoa(i) + "/"
		urls = append(urls, url)
	}
	return urls
}

// Run 抓起页面
func (c *KdlCrawl) Run(url string) error {
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
		var proxy Proxy
		proxy.Ip = s.Find("td[data-title='IP']").First().Text()
		proxy.Port = s.Find("td[data-title='PORT']").First().Text()
		proxy.category = strings.ToLower(s.Find("td[data-title='类型']").First().Text())
		Proxys.Store(proxy.Ip, proxy)
	})
	return nil
}
