package storage

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"

	"pxpool/models"

	bolt "go.etcd.io/bbolt"
)

// Bolt storage
type Bolt struct {
	db   *bolt.DB
	path string
}

// GetBoltStorage 返回storage
func GetBoltStorage(path string) *Bolt {
	db, err := bolt.Open(path+"/pxpool.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}

	return &Bolt{db: db, path: path}
}

// AddOrUpdateProxy 添加或更新
func (b *Bolt) AddOrUpdateProxy(p *models.Proxy) error {
	err := b.db.Update(func(tx *bolt.Tx) error {
		proxy := b.GetProxyByHost(p.Host)
		if proxy == nil { //没有新建
			bProxys, err := tx.CreateBucketIfNotExists([]byte("proxys"))
			if err != nil {
				return err
			}
			_id, err := bProxys.NextSequence()
			if err != nil {
				return err
			}
			proxyID := string(strconv.FormatUint(_id, 10))
			proxyKey := proxyID + "|" + p.Host + ":" + p.Port
			bProxy, err := bProxys.CreateBucketIfNotExists([]byte(proxyKey))
			p.ID = proxyID
			return b.ProxyToBucket(p, bProxy)
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

// GetProxyByHost 更加host查找proxy
func (b *Bolt) GetProxyByHost(host string) *models.Proxy {
	proxy := &models.Proxy{}
	var has bool
	if err := b.db.View(func(tx *bolt.Tx) error {
		bProxys, err := tx.CreateBucketIfNotExists([]byte("proxys"))
		if err != nil {
			return err
		}
		c := bProxys.Cursor()
		search := []byte("|" + host + ":")
		for k, v := c.Seek(search); k != nil && bytes.Contains(k, search); k, v = c.Next() {
			fmt.Printf("key=%s, value=%s\n", k, v)
			bProxy := bProxys.Bucket([]byte(k))
			b.BucketToProxy(bProxy, proxy)
			has = true
			break
		}
		return nil
	}); err != nil {
		return nil
	}
	if has == true {
		return proxy
	}
	return nil
}

// BucketToProxy 赋值
func (b *Bolt) BucketToProxy(bp *bolt.Bucket, p *models.Proxy) error {
	id := string(bp.Get([]byte("ID")))
	log.Println(id)
	p.ID = string(bp.Get([]byte("ID")))
	p.Host = string(bp.Get([]byte("Host")))
	p.Category = string(bp.Get([]byte("Category")))
	p.Port = string(bp.Get([]byte("Port")))
	p.JoinTime = string(bp.Get([]byte("JoinTime")))
	p.VerifyTime = string(bp.Get([]byte("VerifyTime")))
	return nil
}

// ProxyToBucket 导出到
func (b *Bolt) ProxyToBucket(p *models.Proxy, bp *bolt.Bucket) error {
	bp.Put([]byte("ID"), []byte(p.ID))
	bp.Put([]byte("Host"), []byte(p.Host))
	bp.Put([]byte("Category"), []byte(p.Category))
	bp.Put([]byte("Port"), []byte(p.Port))
	bp.Put([]byte("JoinTime"), []byte(p.JoinTime))
	bp.Put([]byte("VerifyTime"), []byte(p.VerifyTime))
	return nil
}

// RandomProxy 随机
func (b *Bolt) RandomProxy() *models.Proxy {
	var l int64
	var proxy = new(models.Proxy)
	if err := b.db.Update(func(tx *bolt.Tx) error {
		bProxys, err := tx.CreateBucketIfNotExists([]byte("proxys"))
		if err != nil {
			return err
		}
		bProxys.ForEach(func(k, v []byte) error {
			l++
			return nil
		})
		n := rand.Int63n(l)
		nString := strconv.FormatInt(n, 10) + "|"
		search := []byte(nString)
		c := bProxys.Cursor()
		for k, v := c.Seek(search); k != nil && bytes.Contains(k, search); k, v = c.Next() {
			fmt.Printf("key=%s, value=%s\n", k, v)
			bProxy := bProxys.Bucket([]byte(k))
			b.BucketToProxy(bProxy, proxy)
			break
		}
		return nil
	}); err != nil {
		log.Println(err)
		return nil
	}
	if proxy.Host != "" {
		return proxy
	}
	return nil
}
