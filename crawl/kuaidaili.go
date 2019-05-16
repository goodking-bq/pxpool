package crawl

import (
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// KdlCrawl 快代理
type KdlCrawl struct {
}

// Start 快代理爬虫
func (c *KdlCrawl) Start() {
	for _, url := range c.GetUrls() {
		go c.crawl(url)
	}
}

// GetUrls 链接
func (c *KdlCrawl) GetUrls() []string {
	var urls []string
	for i := 1; i < 3; i++ {
		urls = append(urls, "https://www.kuaidaili.com/free/inha/"+string(i)+"/")
	}
	return urls
}

func (c *KdlCrawl) crawl(url string) {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	doc.Find("tbody tr").Each(func(i int, s *goquery.Selection) {
		var proxy Proxy
		proxy.Ip = s.Find("td[data-title='IP']").First().Text()
		proxy.Port = s.Find("td[data-title='PORT']").First().Text()
		proxy.category = strings.ToLower(s.Find("td[data-title='类型']").First().Text())
		Proxys.Store(proxy.Ip, proxy)
	})
}
