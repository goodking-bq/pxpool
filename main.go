package main

import (
	"net/http"

	"./crawl"
)

func main() {
	manager := new(crawl.Manager)
	kdl := new(crawl.KdlCrawl)
	manager.Add(kdl)
	http.HandleFunc("/random/", func(w http.ResponseWriter, r *http.Request) {
		p := crawl.Proxys.Random()
		w.Write([]byte(p.URL()))
	})
	//监听3000端口
	http.ListenAndServe(":3000", nil)
}
