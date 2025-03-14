package main

import (
	"context"
	"flag"
	"fmt"

	"mall/service/order/rpc/internal/config"
	"mall/service/order/rpc/internal/server"
	"mall/service/order/rpc/internal/svc"
	"mall/service/order/rpc/internal/task"
	"mall/service/order/rpc/types"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/order.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	// 初始化定时任务
	timeoutTask := task.NewOrderTimeoutTask(context.Background(), ctx)
	timeoutTask.Start()

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		types.RegisterOrderServer(grpcServer, server.NewOrderServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
