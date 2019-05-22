package storage

import "../model"

// Storage 接口
type Storage interface {
	Get() *interface{}
	AddOrUpdateProxy(p *model.Proxy) error
	GetProxyByHost(s string) *model.Proxy
	RandomProxy() *model.Proxy
	Write(p *model.Proxy) error
	Read(p interface{}) *model.Proxy
}

// GetStorage 获取存储
func GetStorage(storageType string) *Storage {
	var storage Storage
	switch storageType {
	case "bolt":
		storage = *GetBoltStorage(".")
	default:
		return nil
	}
	return storage
}
