package storage

import (
	"context"
	"pxpool/models"
	"strings"
)

// Storager 接口
type Storager interface {
	//Get() *interface{}
	AddOrUpdateProxy(p *models.Proxy) error
	GetProxyByHost(host string) *models.Proxy
	RandomProxy() *models.Proxy
	//Write(p *model.Proxy) error
	//Read(p interface{}) *model.Proxy
}

// GetStorage 获取存储
func GetStorage(config *models.Config) *Storager {
	var storage Storager
	switch strings.ToLower(config.StorageType) {
	case "bolt":
		bolt := GetBoltStorage(config.Bolt.DataPath)
		storage = bolt
		break
	default:
		return nil
	}
	return &storage
}

// StartStorage 保存扫描到的代理
func StartStorage(ctx context.Context, storage *Storager, dataChan chan *models.Proxy) {
	for {
		select {
		case proxy := <-dataChan:
			(*storage).AddOrUpdateProxy(proxy)
		case <-ctx.Done():
			close(dataChan)
		}
	}
}
