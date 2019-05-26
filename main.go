package main

import (
	"context"
	"os"
	"pxpool/scanner"

	"pxpool/crawler"
	"pxpool/models"
	"pxpool/storage"
	"pxpool/web"

	"github.com/urfave/cli"
)

// WebAction web 动作
func WebAction(c *cli.Context) error {
	config := models.DefaultConfig()
	config.UnmarshalCtx(c)
	ctx, cancal := context.WithCancel(context.Background())
	defer cancal()
	dataChan := make(chan *models.Proxy)
	scanner := scanner.NewScanner(config)
	go scanner.Scan(ctx, config, dataChan)
	storager := storage.GetStorage(config)
	go storage.StartStorage(ctx, storager, dataChan)
	api := web.DefaultAPI(*storager)
	api.Run(ctx, config)
	return nil
}

// ScanAction scan 动作
func ScanAction(c *cli.Context) error {
	config := models.ConfigFromCtx(c)
	ctx, cancal := context.WithCancel(context.Background())
	defer cancal()
	dataChan := make(chan *models.Proxy)
	scanner := scanner.NewScanner(config)
	go scanner.Scan(ctx, config, dataChan)
	storager := storage.GetStorage(config)
	storage.StartStorage(ctx, storager, dataChan)
	// exit := make(chan os.Signal)
	// signal.Notify(exit, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGKILL)
	// signal := <-exit
	// cancal()
	// log.Println(signal)
	return nil
}

func main() {
	pxApp := cli.NewApp()
	pxApp.Name = "代理扫描工具"
	pxApp.Version = "0.1"
	pxApp.Usage = "代理站全功能"
	//pxApp.UsageText = "什么什么"

	pxApp.Commands = []cli.Command{
		{
			Name:  "web",
			Usage: "启动webapi",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "bind",
					Usage: "绑定ip",
				},
				cli.StringFlag{
					Name:  "port",
					Usage: "端口",
				},
			},
			Action: WebAction,
		},
		{
			Name:  "crawl",
			Usage: "启动爬虫进程",
			Flags: []cli.Flag{},
			Action: func(c *cli.Context) {
				config := models.ConfigFromCtx(c)
				ctx, cancal := context.WithCancel(context.Background())
				defer cancal()
				storager := storage.GetStorage(config)
				DataChan := make(chan *models.Proxy)
				cManager := crawler.NewDefaultCrawl(storager, DataChan)
				cManager.Start()
				go storage.StartStorage(ctx, storager, DataChan)
				api := web.DefaultAPI(*storager)
				api.Run(ctx, config)
				cManager.ExitSignal <- true
			},
		},
		{
			Name:  "scan",
			Usage: "启动代理扫描",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "scanfile",
					Usage: "",
					Value: "scan.txt",
				},
				cli.StringFlag{
					Name:  "cidr",
					Usage: "ru 172.0.0.1/24",
				},
				cli.IntFlag{
					Name:  "concurrency,C",
					Usage: "扫描进程数",
					Value: 100,
				},
			},
			Action: ScanAction,
		},
		{
			Name:  "all",
			Usage: "启动所有",
		},
	}
	pxApp.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config,c",
			Usage: "配置文件路径",
		},
		cli.StringFlag{
			Name:  "datapath,d",
			Usage: "数据文件目录",
		},
	}
	pxApp.Run(os.Args)
}
