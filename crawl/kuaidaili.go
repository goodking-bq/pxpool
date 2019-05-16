package crawl

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// KdlCrawl 快代理
type KdlCrawl struct {
}

// Start 快代理爬虫
func (c *KdlCrawl) Start() {
	for _, url := range c.GetUrls() {
		go c.Run(url)
	}
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
func (c *KdlCrawl) Run(url string) {
	fmt.Println(url)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, strings.NewReader(""))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows; U; Windows NT 6.1; en-US) AppleWebKit/532.5 (KHTML, like Gecko) Chrome/4.0.249.0 Safari/532.5")
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", resp.StatusCode, resp.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(resp.Body)
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
	fmt.Println("done ...")
}
