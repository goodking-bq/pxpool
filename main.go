package main

import (
	"os"

	"github.com/urfave/cli"
)

func main() {
	pxApp := cli.NewApp()
	pxApp.Name = "代理扫描工具"
	pxApp.Version = "0.1"
	pxApp.Usage = "代理站全功能"
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
		},
		{
			Name:  "crawl",
			Usage: "启动爬虫进程",
			Flags: []cli.Flag{},
		},
		{
			Name:  "scanner",
			Usage: "启动代理扫描",
		},
		{
			Name:  "all",
			Usage: "启动所有",
		},
	}
	pxApp.Run(os.Args)

}
