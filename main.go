package main

import (
	"./crawl"
	"./scanner"
	"./web"
)

func main() {
	crawlManager := crawl.NewDefaultManager()
	crawlManager.Start()
	scan := scanner.Scanner{}
	scan.ScanCidr([]byte("127.0.0.0/30"))
	api := web.NewDefaultAPI()
	api.Run("", 3000)
}
