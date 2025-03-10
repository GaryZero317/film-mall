package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

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

	server := rest.MustNewServer(c.RestConf, rest.WithCustomCors(nil, func(w http.ResponseWriter) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		w.Header().Set("Access-Control-Expose-Headers", "Content-Length, Content-Type, Access-Control-Allow-Origin, Access-Control-Allow-Headers")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
	}, "*"))

	defer server.Stop()

	// 添加静态文件服务 - 使用与product服务相同的简单方式
	uploadDir := "D:/graduation/FilmMall/service/film/api/uploads"
	logx.Infof("静态文件服务设置: 目录=%s", uploadDir)

	// 确保目录存在
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		logx.Errorf("创建上传目录失败: %v", err)
	}

	// 直接注册静态文件路由, 使用简单的http.FileServer
	server.AddRoute(rest.Route{
		Method: http.MethodGet,
		Path:   "/uploads/:file",
		Handler: http.StripPrefix("/uploads/",
			http.FileServer(http.Dir(uploadDir))).ServeHTTP,
	})

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
	logx.Infof("胶片冲洗服务已启动，监听地址: %s:%d", c.Host, c.Port)
	server.Start()
}

// 获取文件大小（保留这个辅助函数，以便在其他地方可能会用到）
func getFileSize(dir string, file os.DirEntry) int64 {
	if !file.IsDir() {
		if info, err := file.Info(); err == nil {
			return info.Size()
		}
	}
	return 0
}
