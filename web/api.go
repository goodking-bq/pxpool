package web

import (
	"fmt"
	"net/http"

	"../model"
)

// API somethiing
type API struct {
	bind string
	port int
}

// NewDefaultAPI 默认
func NewDefaultAPI() *API {
	return &API{bind: "", port: 3000}
}

// Run 启动api
func (api *API) Run(bind string, port int) {
	http.HandleFunc("/random/", func(w http.ResponseWriter, r *http.Request) {
		p, err := model.GetProxyStory().Random()
		if err != nil {
			w.Write([]byte("err"))
		} else {
			w.Write([]byte(p.URL()))
		}

	})
	//监听3000端口
	http.ListenAndServe(fmt.Sprintf("%s:%d", bind, port), nil)
}
