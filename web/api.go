package web

import (
	"fmt"
	"net/http"

	"../storage"
)

// API somethiing
type API struct {
	bind    string
	port    int
	storage storage.Storager
}

// NewDefaultAPI 默认
func NewDefaultAPI(storage storage.Storager) *API {
	return &API{bind: "", port: 3000, storage: storage}
}

// Run 启动api
func (api *API) Run(bind string, port int) {
	http.HandleFunc("/random/", func(w http.ResponseWriter, r *http.Request) {
		p := api.storage.RandomProxy()
		if p == nil {
			w.Write([]byte("err"))
		} else {
			w.Write([]byte(p.URL()))
		}

	})
	//监听3000端口
	http.ListenAndServe(fmt.Sprintf("%s:%d", bind, port), nil)
}
