package main

import (
	"flag"
	"fmt"
	"net/http"

	"mall/service/film/api/internal/config"
	"mall/service/film/api/internal/handler"
	"mall/service/film/api/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/film.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	// 打印JWT配置信息
	logx.Infof("JWT Auth配置: AccessSecret长度=%d, AccessExpire=%d", len(c.Auth.AccessSecret), c.Auth.AccessExpire)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	// 添加一个自定义中间件，用于记录请求头信息
	server.Use(func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// 打印请求信息
			logx.Infof("收到请求: %s %s", r.Method, r.URL.Path)

			// 记录Authorization头并分析JWT
			auth := r.Header.Get("Authorization")
			if auth != "" {
				if len(auth) > 15 {
					// 只打印前15个字符，避免泄露完整令牌
					logx.Infof("收到Authorization头: %s...（已截断）", auth[:15])

					// 分析JWT令牌
					svc.ParseJwtToken(auth, c.Auth.AccessSecret)
				} else {
					logx.Infof("收到Authorization头格式可能不正确: %s", auth)
				}
			} else {
				logx.Error("请求中缺少Authorization头")
			}

			next(w, r)
		}
	})

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	logx.Infof("Film服务已启动，监听地址: %s:%d", c.Host, c.Port)
	server.Start()
}
