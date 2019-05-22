package model

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	bolt "go.etcd.io/bbolt"
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
		dialer, err := nproxy.SOCKS5("tcp", p.Host+":"+p.Port, nil, nproxy.Direct)
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

// Save 保存到db
func (p *Proxy) Save() error {
	var proxyID string
	if err := Storage.db.Update(func(tx *bolt.Tx) error {
		bProxys, err := tx.CreateBucketIfNotExists([]byte("proxys"))
		if err != nil {
			return err
		}
		if p.ID == "" {
			_id, err := bProxys.NextSequence()
			proxyID = strconv.FormatUint(_id, 10)
			if err != nil {
				return err
			}
			p.ID = proxyID
		} else {
			proxyID = p.ID
		}
		bProxy, err := bProxys.CreateBucketIfNotExists([]byte(proxyID))
		if err != nil {
			return err
		}
		return p.MarshalBucket(bProxy)
	}); err != nil {
		return err
	}

	log.Printf("save %s success", p.Host)
	return nil
}

func (p *Proxy) MarshalBucket(bucket *bolt.Bucket) error {
	bucket.Put([]byte("id"), []byte(p.ID))
	bucket.Put([]byte("host"), []byte(p.Host))
	bucket.Put([]byte("post"), []byte(p.Port))
	bucket.Put([]byte("category"), []byte(p.Category))
	bucket.Put([]byte("jointime"), []byte(p.JoinTime))
	bucket.Put([]byte("verifytime"), []byte(p.VerifyTime))
	return nil
}

func (p *Proxy) UnMarshalBucket(bucket *bolt.Bucket) error {
	str := string(bucket.Get([]byte("id")))
	p.ID = str
	p.Host = string(bucket.Get([]byte("host")))
	p.Port = string(bucket.Get([]byte("post")))
	p.Category = string(bucket.Get([]byte("category")))
	p.JoinTime = string(bucket.Get([]byte("jointime")))
	p.VerifyTime = string(bucket.Get([]byte("verifytime")))
	return nil
}

// ProxyStory 代理管理器
type ProxyStory struct {
}

//GetProxyStory huoqu
func GetProxyStory() *ProxyStory {
	return &ProxyStory{}
}

// Random 随机获取
func (p *ProxyStory) Random() (*Proxy, error) {
	var l int64
	var proxy = new(Proxy)
	if err := Storage.db.View(func(tx *bolt.Tx) error {
		bProxys := tx.Bucket([]byte("proxys"))
		bProxys.ForEach(func(k, v []byte) error {
			l++
			return nil
		})
		n := rand.Int63n(l)
		nString := strconv.FormatInt(n, 10)
		nbyte := []byte(nString)
		bProxy := bProxys.Bucket(nbyte)
		return proxy.UnMarshalBucket(bProxy)
	}); err != nil {
		log.Println(err)
		return nil, err
	}
	return proxy, nil
}
