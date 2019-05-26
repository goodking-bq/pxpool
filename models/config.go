package models

import (
	"io/ioutil"

	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"
)

type bolt struct {
	DataPath string
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
		StorageType: "",
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
	var (
		configFile = ctx.Parent().String("config")
		dataPath   = ctx.Parent().String("datapath")
	)
	if configFile != "" {
		conf, err := ioutil.ReadFile(configFile)
		if err != nil {
			return err
		}
		err = yaml.Unmarshal(conf, c)
		if err != nil {
			print(err)
			return err
		}
	}
	if dataPath != "" {
		c.Bolt.DataPath = dataPath
	}
	return nil
}
