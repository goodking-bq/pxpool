package web

import (
	"pxpool/models"
	"pxpool/storage"
	"strconv"

	"github.com/sirupsen/logrus"

	"github.com/kataras/iris"
)

type WebApp struct {
	App     *iris.Application
	Storage storage.Storager
	secret  string
	logger  *logrus.Logger
}

func NewApp(storage storage.Storager, secret string, logger *logrus.Logger) *WebApp {
	return &WebApp{
		App:     iris.New(),
		Storage: storage,
		secret:  secret,
		logger:  logger,
	}
}

func (app *WebApp) Run() {
	apis := app.App.Party("/api")
	apis.Get("/random", func(ctx iris.Context) {
		px := app.Storage.RandomProxy()
		if px != nil {
			ctx.Write([]byte(px.URL()))
		} else {
			ctx.Write([]byte("error: no proxy"))
		}

	})
	apis.Get("/counter", func(ctx iris.Context) {
		p := app.Storage.GetProxyCounter()
		ctx.Write([]byte(strconv.FormatInt(p, 10)))
	})
	apis.Post("/proxy", func(ctx iris.Context) {
		host := ctx.PostValue("ip")
		port := ctx.PostValue("port")
		category := ctx.PostValue("category")
		proxy := models.NewProxy(host, port, category)
		app.Storage.AddOrUpdateProxy(proxy)
		app.logger.WithField("proxy", proxy.URL()).Infoln("保存成功")
	})

	home := app.App.Party("/")
	home.Get("search", func(ctx iris.Context) {
		ip := ctx.URLParam("ip")
		p := app.Storage.GetProxysByHost(ip)
		if len(p) == 0 {
			ctx.Write([]byte("proxy not exist"))
		} else {
			for _, px := range p {
				ctx.Write([]byte(px.URL() + "\n"))
			}
		}
	})

	app.App.Run(iris.Addr(":3000"))
}
