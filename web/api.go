package web

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"pxpool/models"
	"pxpool/storage"
)

// API somethiing
type API struct {
	bind    string
	port    int
	storage storage.Storager
}

// DefaultAPI 默认
func DefaultAPI(storage storage.Storager) *API {
	return &API{bind: "", port: 3000, storage: storage}
}

// Run 启动api
func (api *API) Run(ctx context.Context, config *models.Config) {
	http.HandleFunc("/random/", func(w http.ResponseWriter, r *http.Request) {
		p := api.storage.RandomProxy()
		if p == nil {
			w.Write([]byte("err"))
		} else {
			w.Write([]byte(p.URL()))
		}

	})

	http.HandleFunc("/counter/", func(w http.ResponseWriter, r *http.Request) {
		p := api.storage.GetProxyCounter()
		w.Write([]byte(strconv.FormatInt(p, 10)))
	})
	http.HandleFunc("/search/", func(w http.ResponseWriter, r *http.Request) {
		p := api.storage.GetProxyByHost("27.203.247.35")
		w.Write([]byte(p.ID))
	})
	//监听3000端口
	http.ListenAndServe(fmt.Sprintf("%s:%d", config.Web.Bind, config.Web.Port), nil)
}
