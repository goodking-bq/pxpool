package storage

import (
	"log"
	"strconv"

	"../model"
	bolt "go.etcd.io/bbolt"
)

// Bolt storage
type Bolt struct {
	db   *bolt.DB
	path string
}

// GetBoltStorage 返回storage
func GetBoltStorage(path string) *Bolt {
	db, err := bolt.Open(path+"/pxpool.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	return &Bolt{db: db, path: path}
}

// AddOrUpdateProxy 添加或更新
func (b *Bolt) AddOrUpdateProxy(p *model.Proxy) error {
	var proxyID string
	err := b.db.Update(func(tx *bolt.Tx) error {
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
		return nil
	})
	return nil

}

// GetProxyByHost 更加host查找proxy
func (b *Bolt) GetProxyByHost(s string) *model.Proxy {
	return nil
}
