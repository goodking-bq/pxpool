package storage

import (
	"../model"
)

// Storager 接口
type Storager interface {
	//Get() *interface{}
	AddOrUpdateProxy(p *model.Proxy) error
	GetProxyByHost(host string) *model.Proxy
	RandomProxy() *model.Proxy
	//Write(p *model.Proxy) error
	//Read(p interface{}) *model.Proxy
}

// GetStorage 获取存储
func GetStorage(storageType string) *Storager {
	var storage Storager
	switch storageType {
	case "bolt":
		bolt := GetBoltStorage(".")
		storage = bolt
		break
	default:
		return nil
	}
	return &storage
}
