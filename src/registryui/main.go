package main

import (
	"flag"
	"log"
	"registryui/client"
	"registryui/view"

	"github.com/labstack/echo/v4"
)

var (
	basepath = "/registry"
	debug    = true
	registry = "192.168.1.254:5000"
	access   = "https://registry:5000"
	listen   = "0.0.0.0:8080"
)

func init() {
	flag.StringVar(&basepath, "basepath", basepath, "Base path of Docker Registry UI")
	flag.BoolVar(&debug, "debug", debug, "Debug mode. Affects only templates")
	flag.StringVar(&registry, "registry", registry, "Registry URL with schema and port for view")
	flag.StringVar(&access, "access", access, "Registry access address")
	flag.StringVar(&listen, "listen", listen, "Listen address")
}

func main() {
	flag.Parse()

	// 初始化数据层，用户获取registry数据
	r, err := client.NewRegistry(access)
	if err != nil {
		log.Fatal(err)
	}
	model := view.NewModel(basepath, r)

	// 初始化模板引擎
	e := echo.New()
	e.Renderer = view.SetupRender(debug, registry, basepath)

	e.Static(basepath+"/static", "resources/static")
	e.File("/favicon.ico", "resources/static/favicon.ico")
	e.GET("", model.Redrict)
	e.GET(basepath, model.Repositories)
	e.GET(basepath+"/:repo", model.Tags)
	e.GET(basepath+"/:repo/:tag", model.Manifest)

	// 启动服务
	log.Fatal(e.Start(listen))
}
