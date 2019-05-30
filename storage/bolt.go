package storage

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"sync"
	"time"

	"pxpool/models"

	bolt "go.etcd.io/bbolt"
)

// Bolt storage
type Bolt struct {
	db   *bolt.DB
	path string
}

var db *Bolt
var once sync.Once

// GetBoltStorage 返回storage 单例模式
func GetBoltStorage(path string) *Bolt {
	once.Do(func() {
		thisdb, err := bolt.Open(path+"/pxpool.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
		if err != nil {
			log.Fatal(err)
		}
		db = &Bolt{db: thisdb, path: path}
	})

	return db
}

// AddOrUpdateProxy 添加或更新
func (b *Bolt) AddOrUpdateProxy(p *models.Proxy) error {
	isNew := false
	err := b.db.Batch(func(tx *bolt.Tx) error {
		proxy := b.GetProxyByproxy(p)
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
			err = b.ProxyToBucket(p, bProxy)
			isNew = true
			return err

		}
		return nil
	})
	if err != nil {
		return err
	}
	if isNew == true {
		b.IncProxyCounter() // 不能放update里面
	}
	return nil
}

// GetProxysByHost 根据host查找proxy
func (b *Bolt) GetProxysByHost(host string) []*models.Proxy {
	var proxys []*models.Proxy
	if err := b.db.View(func(tx *bolt.Tx) error {
		bProxys := tx.Bucket([]byte("proxys"))
		if bProxys == nil {
			return errors.New("no proxys bucket")
		}
		c := bProxys.Cursor()
		searchstr := "|" + host + ":"
		print(searchstr)
		search := []byte(searchstr)

		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			if bytes.Contains(k, search) {
				proxy := &models.Proxy{}
				fmt.Printf("key=%s\n", k)
				bProxy := bProxys.Bucket([]byte(k))
				b.BucketToProxy(bProxy, proxy)
				proxys = append(proxys, proxy)
			}

		}
		return nil
	}); err != nil {
		return nil
	}

	return proxys

}

// GetProxyByproxy xxx
func (b *Bolt) GetProxyByproxy(p *models.Proxy) *models.Proxy {
	if err := b.db.View(func(tx *bolt.Tx) error {
		bProxys := tx.Bucket([]byte("proxys"))
		if bProxys == nil {
			return errors.New("no proxys bucket")
		}
		c := bProxys.Cursor()
		searchstr := "|" + p.Host + ":" + p.Port
		print(searchstr)
		search := []byte(searchstr)

		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			if bytes.Contains(k, search) {
				fmt.Printf("key=%s\n", k)
				bProxy := bProxys.Bucket([]byte(k))
				b.BucketToProxy(bProxy, p)
			}

		}
		return nil
	}); err != nil {
		return nil
	}
	if p.ID != "" {
		return p
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
	if err := b.db.View(func(tx *bolt.Tx) error {
		bProxys := tx.Bucket([]byte("proxys"))
		if bProxys == nil {
			return errors.New("no proxys bucket")
		}
		bProxys.ForEach(func(k, v []byte) error {
			l++
			return nil
		})
		if l == 0 {
			return nil
		}
		n := rand.Int63n(l)
		n++
		nString := strconv.FormatInt(n, 10) + "|"
		search := []byte(nString)
		c := bProxys.Cursor()
		for k, v := c.Seek(search); k != nil && bytes.HasPrefix(k, search); k, v = c.Next() {
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

// SetProxyCounter proxy 计数器设置
func (b *Bolt) SetProxyCounter(n int64) error {
	err := b.db.Update(func(tx *bolt.Tx) error {
		counter, err := tx.CreateBucketIfNotExists([]byte("proxycounter"))
		if err != nil {
			return err
		}
		counter.Put([]byte("count"), []byte(strconv.FormatInt(n, 10)))
		counter.Put([]byte("update"), []byte(time.Now().String()))
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

// GetProxyCounter 获取proxy个数
func (b *Bolt) GetProxyCounter() int64 {
	var count int64
	b.db.View(func(tx *bolt.Tx) error {
		counter := tx.Bucket([]byte("proxycounter"))
		if counter != nil {
			buf := counter.Get([]byte("count"))
			count1, err := strconv.ParseInt(string(buf), 10, 64)
			if err != nil {
				return err
			}
			count = count1
		}

		return nil
	})
	return count
}

// IncProxyCounter counter +1
func (b *Bolt) IncProxyCounter() error {
	err := b.db.Batch(func(tx *bolt.Tx) error {
		counter, err := tx.CreateBucketIfNotExists([]byte("proxycounter"))
		if err != nil {
			return err
		}
		buf := counter.Get([]byte("count"))
		count, err := strconv.ParseInt(string(buf), 10, 64)
		if err != nil {
			count = 0
		}
		count++
		log.Println(count)
		counter.Put([]byte("count"), []byte(strconv.FormatInt(count, 10)))
		counter.Put([]byte("update"), []byte(time.Now().String()))
		return nil
	})
	return err
}
