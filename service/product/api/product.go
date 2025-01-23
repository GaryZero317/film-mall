package main

import (
	"flag"
	"fmt"
	"net/http"

	"mall/service/product/api/internal/config"
	"mall/service/product/api/internal/handler"
	"mall/service/product/api/internal/svc"
	"mall/service/product/model"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var configFile = flag.String("f", "etc/product.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	// 初始化 GORM
	db, err := gorm.Open(mysql.Open(c.Mysql.DataSource), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("连接数据库失败: %v", err))
	}

	// 执行自动迁移
	if err := model.AutoMigrate(db); err != nil {
		panic(fmt.Sprintf("数据库迁移失败: %v", err))
	}

	server := rest.MustNewServer(c.RestConf, rest.WithCustomCors(nil, func(w http.ResponseWriter) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		w.Header().Set("Access-Control-Expose-Headers", "Content-Length, Content-Type, Access-Control-Allow-Origin, Access-Control-Allow-Headers")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
	}, "*"))

	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	// 注册静态文件处理
	server.AddRoute(rest.Route{
		Method: http.MethodGet,
		Path:   "/uploads/:file",
		Handler: http.StripPrefix("/uploads/",
			http.FileServer(http.Dir("D:/graduation/FilmMall/service/product/api/uploads"))).ServeHTTP,
	})

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
