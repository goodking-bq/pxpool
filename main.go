package main

import (
	"./crawl"
	"./web"
)

func main() {
	crawlManager := crawl.NewDefaultManager()
	crawlManager.Start()
	api := web.NewDefaultAPI()
	api.Run("", 3000)
}
