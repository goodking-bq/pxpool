package models

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	nproxy "golang.org/x/net/proxy"
)

// Proxy 代理
type Proxy struct {
	ID         string
	Host       string
	Port       string
	Category   string
	JoinTime   string
	VerifyTime string
}

// ProxyChan Proxy channel
var ProxyChan = make(chan *Proxy)

// CheckProxy 检查代理是否可用
func CheckProxy(p *Proxy) bool {
	return p.check()
}

// NewProxy 创建新的
func NewProxy(host, port, category string) *Proxy {
	return &Proxy{
		Host:       host,
		Port:       port,
		Category:   category,
		JoinTime:   time.Now().String(),
		VerifyTime: time.Now().String(),
	}
}

// URL 获取代理的地址
func (p *Proxy) URL() string {
	return p.Category + "://" + p.Host + ":" + p.Port
}

func (p *Proxy) Key() string {
	return p.Host + ":" + p.Port
}

func (p *Proxy) check() bool {
	log.Printf("开始检查 ..,%s", p.URL())
	netTransport := &http.Transport{
		TLSHandshakeTimeout:   1 * time.Second,
		ResponseHeaderTimeout: 2 * time.Second,
		MaxIdleConnsPerHost:   1,
	}
	if strings.ToLower(p.Category) == "http" || strings.ToLower(p.Category) == "https" {
		proxy, _ := url.Parse(p.URL())
		netTransport.Proxy = http.ProxyURL(proxy)
		netTransport.Dial = func(netw, addr string) (net.Conn, error) {
			c, err := net.DialTimeout(netw, addr, time.Second*time.Duration(2))
			if err != nil {
				return nil, err
			}
			return c, nil
		}
	} else {
		dialer, err := nproxy.SOCKS5("tcp", p.Host+":"+p.Port, nil, &net.Dialer{Timeout: 2 * time.Second})
		if err != nil {
			fmt.Fprintln(os.Stderr, "can't connect to the proxy:", err)
			return false
		}
		// setup a http client
		netTransport.Dial = dialer.Dial
	}

	client := &http.Client{
		Timeout:   time.Second * 2,
		Transport: netTransport}
	req, err := http.NewRequest("GET", "http://www.ip138.com/", nil)
	defer client.CloseIdleConnections()
	if err != nil {
		return true
	}
	req.Close = true
	req.Header.Add("Accept-Encoding", "identity")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows; U; Windows NT 6.1; en-US) AppleWebKit/532.5 (KHTML, like Gecko) Chrome/4.0.249.0 Safari/532.5")
	log.Printf("使用代理链接中 ..,%s", p.URL())
	resp, err := client.Do(req)
	log.Printf("使用代理链接中 ..,%s -- 完成", p.URL())
	if err != nil {
		log.Println(err)
		return false
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return false
	}
	return true
}

// ProxyStory 代理管理器
type ProxyStory struct {
}

//GetProxyStory huoqu
func GetProxyStory() *ProxyStory {
	return &ProxyStory{}
}
