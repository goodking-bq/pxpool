package storage

import "../model"

type Storage interface {
	Get() *interface{}
	SaveProxy(p *model.Proxy) error
	GetProxy(s string) *model.Proxy
	RandomProxy() *model.Proxy
}
