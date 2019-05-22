package main

import (
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
				
			}
		},
	}
	pxApp.Run(os.Args)

}
