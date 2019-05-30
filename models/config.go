package models

import (
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/spf13/viper"
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
	Ports          []int
	PortString     string
}

type logger struct {
	File  string
	Level string
}

type crawl struct {
	Ticker int
}

// Config 配置
type Config struct {
	StorageType string
	Bolt        bolt
	Web         web
	Scanner     scanner
	Crawl       crawl
	Log         logger
	Url         string
	Secret      string
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

// UnmarshalViper 从viper中加载配置
func (c *Config) UnmarshalViper() {
	c.StorageType = viper.GetString("storagetype")
	c.Log.File = viper.GetString("log.file")
	c.Log.Level = viper.GetString("log.level")
	c.Url = viper.GetString("url")
	c.Secret = viper.GetString("secret")
	c.Crawl.Ticker = viper.GetInt("crawl.ticker")
	c.Scanner.Cidr = viper.GetString("scanner.cidr")
	c.Scanner.File = viper.GetString("scanner.file")
	c.Scanner.PortString = viper.GetString("scanner.portstring")
	if c.Scanner.PortString != "" {
		var ports []int
		for _, port := range strings.Split(c.Scanner.PortString, ",") {
			i, _ := strconv.Atoi(port)
			ports = append(ports, i)
		}
		c.Scanner.Ports = ports
	} else {
		for _, p := range strings.Split(viper.GetStringSlice("scanner.ports")[0], ",") {
			i, _ := strconv.Atoi(p)
			c.Scanner.Ports = append(c.Scanner.Ports, i)
		}
	}

	c.Scanner.MaxConcurrency = viper.GetInt("scanner.maxconcurrency")
	c.Web.Bind = viper.GetString("web.bind")
	c.Web.Port = viper.GetInt("web.port")
	c.Bolt.DataPath = viper.GetString("bolt.datapath")
}
