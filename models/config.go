package models

import (
	"github.com/urfave/cli"
)

type bolt struct {
	DataPath string `datapath: 数据文件目录`
}
type web struct {
	Bind string
	Port int
}
type scanner struct {
	File           string
	Cidr           string
	MaxConcurrency int
}

// Config 配置
type Config struct {
	StorageType string
	Bolt        bolt
	Web         web
	Scanner     scanner
}

// DefaultConfig 默认配置
func DefaultConfig() *Config {
	return &Config{
		StorageType: "bolt",
		Bolt: bolt{
			DataPath: ".",
		},
		Web: web{
			Bind: "0.0.0.0",
			Port: 3000,
		},
		Scanner: scanner{},
	}
}

//ConfigFromCtx 配置文件
func ConfigFromCtx(ctx *cli.Context) *Config {
	config := DefaultConfig()
	config.UnmarshalCtx(ctx)
	return config
}

// UnmarshalCtx 从cli.Context 获取配置
func (c *Config) UnmarshalCtx(ctx *cli.Context) error {
	c.Scanner.Cidr = ctx.String("cidr")
	c.Scanner.File = ctx.String("scanfile")
	c.Scanner.MaxConcurrency = ctx.Int("concurrency")
	return nil
}
