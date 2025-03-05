package main

import (
	"flag"
	"fmt"
	"net/http"

	"mall/service/community/api/internal/config"
	"mall/service/community/api/internal/handler"
	"mall/service/community/api/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/community.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(rest.RestConf{
		MaxBytes: 1024 * 1024 * 1024, // 设置最大请求体大小为1GB
		Host:     c.Host,
		Port:     c.Port,
		Timeout:  c.Timeout,
	}, rest.WithCustomCors(nil, func(w http.ResponseWriter) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		w.Header().Set("Access-Control-Expose-Headers", "Content-Length, Content-Type, Access-Control-Allow-Origin, Access-Control-Allow-Headers")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		logx.Info("已设置CORS响应头")
	}, "*"))

	ctx := svc.NewServiceContext(c)
	defer server.Stop()

	// 添加静态文件服务
	fileServer := http.FileServer(http.Dir("api/uploads"))
	server.AddRoute(rest.Route{
		Method: http.MethodGet,
		Path:   "/uploads/works/:file",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.StripPrefix("/uploads/", fileServer).ServeHTTP(w, r)
		}),
	})

	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
