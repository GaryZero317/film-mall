package main

import (
	"flag"
	"net/http"

	"mall/service/user/api/internal/config"
	"mall/service/user/api/internal/handler"
	"mall/service/user/api/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/user.yaml", "the config file")

func main() {
	flag.Parse()

	logx.Info("开始加载配置文件")
	var c config.Config
	conf.MustLoad(*configFile, &c)
	logx.Infof("配置文件加载成功，配置内容: %+v", c)

	logx.Info("创建服务上下文")
	ctx := svc.NewServiceContext(c)
	logx.Info("服务上下文创建成功")

	logx.Info("创建HTTP服务器")
	server := rest.MustNewServer(c.RestConf, rest.WithCustomCors(nil, func(w http.ResponseWriter) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		w.Header().Set("Access-Control-Expose-Headers", "Content-Length, Content-Type, Access-Control-Allow-Origin, Access-Control-Allow-Headers")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		logx.Info("已设置CORS响应头")
	}, "*"))
	defer server.Stop()
	logx.Info("HTTP服务器创建成功")

	logx.Info("注册路由处理器")
	handler.RegisterHandlers(server, ctx)
	logx.Info("路由处理器注册成功")

	logx.Infof("启动服务器在 %s:%d...", c.Host, c.Port)
	logx.Info("JWT认证配置: ")
	logx.Infof("  用户认证密钥: %s, 有效期: %d 秒", c.Auth.AccessSecret, c.Auth.AccessExpire)
	logx.Infof("  管理员认证密钥: %s, 有效期: %d 秒", c.AdminAuth.AccessSecret, c.AdminAuth.AccessExpire)

	logx.Info("开始监听HTTP请求")
	server.Start()
}
