package model

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/proxy"
)

// Proxy 代理
type Proxy struct {
	Ip         string
	Port       string
	Category   string
	joinTime   string
	verifyTime string
}

// URL 获取代理的地址
func (p *Proxy) URL() string {
	return p.Category + "://" + p.Ip + ":" + p.Port
}

// CheckProxy 检查代理是否可用
func CheckProxy(p *Proxy) bool {
	return p.check()
}
func (p *Proxy) check() bool {
	netTransport := &http.Transport{}
	if strings.ToLower(p.Category) == "http" || strings.ToLower(p.Category) == "https" {
		proxy, _ := url.Parse(p.URL())
		netTransport.Proxy = http.ProxyURL(proxy)
		netTransport.Dial = func(netw, addr string) (net.Conn, error) {
			c, err := net.DialTimeout(netw, addr, time.Second*time.Duration(10))
			if err != nil {
				return nil, err
			}
			return c, nil
		}
		//Proxy: http.ProxyFromEnvironment,

		netTransport.MaxIdleConnsPerHost = 10                               //每个host最大空闲连接
		netTransport.ResponseHeaderTimeout = time.Second * time.Duration(5) //数据收发5秒超时

	} else {
		dialer, err := proxy.SOCKS5("tcp", p.Ip+":"+p.Port, nil, proxy.Direct)
		if err != nil {
			fmt.Fprintln(os.Stderr, "can't connect to the proxy:", err)
			os.Exit(1)
		}
		// setup a http client
		netTransport.Dial = dialer.Dial
	}

	client := &http.Client{
		Timeout:   time.Second * 10,
		Transport: netTransport}
	req, err := http.NewRequest("GET", "http://www.ip138.com/", nil)
	if err != nil {
		return true
	}
	req.Close = true
	req.Header.Add("Accept-Encoding", "identity")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows; U; Windows NT 6.1; en-US) AppleWebKit/532.5 (KHTML, like Gecko) Chrome/4.0.249.0 Safari/532.5")
	resp, err := client.Do(req)
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

type proxysMap struct {
	sync.Map
}

// Proxys 所有的代理
var Proxys proxysMap

func (p *proxysMap) Random() (Proxy, error) {
	var ips []string
	Proxys.Range(func(k, _ interface{}) bool {
		ips = append(ips, k.(string))
		return true
	})
	l := len(ips)
	if l == 0 {
		return Proxy{}, errors.New("没有缓存代理")
	}
	n := rand.Intn(l)
	_p, _ := Proxys.Load(ips[n])
	return _p.(Proxy), nil
}
func (p *proxysMap) Check() {
	log.Println("开始检查 。。。")
	Proxys.Range(func(k, v interface{}) bool {
		px := v.(Proxy)
		isActive := px.check()
		if !isActive {
			Proxys.Delete(k)
		}
		return true
	})

}
