package model

// Manager 爬虫管理器
type Manager interface {
	Add(i interface{}) error
	Start()
	StartTicker() chan bool
	StartAndTicker() chan bool
}
