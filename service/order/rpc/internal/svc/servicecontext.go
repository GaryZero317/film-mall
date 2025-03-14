package svc

import (
	"mall/service/order/model"
	"mall/service/order/rpc/internal/config"
	"mall/service/pay/rpc/pay"
	"mall/service/product/rpc/pb/product"

	"github.com/go-redis/redis/v8"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config          config.Config
	OrderModel      model.OrderModel
	OrderItemModel  model.OrderItemModel
	Redis           *redis.Client
	ProductRpc      product.ProductClient
	PayRpc          pay.Pay
	OrderLockExpiry int
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Mysql.DataSource)

	// 初始化Redis客户端
	rdb := redis.NewClient(&redis.Options{
		Addr:     c.OrderRedis.Host,
		Password: c.OrderRedis.Pass,
		DB:       0,
	})

	ctx := &ServiceContext{
		Config:          c,
		OrderModel:      model.NewOrderModel(conn, c.CacheRedis),
		OrderItemModel:  model.NewOrderItemModel(conn, c.CacheRedis),
		Redis:           rdb,
		ProductRpc:      product.NewProductClient(zrpc.MustNewClient(c.ProductRpc).Conn()),
		OrderLockExpiry: c.OrderLockExpiry,
	}

	// 可选地连接支付服务
	if c.PayRpc.Etcd.Hosts != nil && len(c.PayRpc.Etcd.Hosts) > 0 {
		// 尝试连接支付服务，如果失败不影响订单服务启动
		payRpcClient, err := zrpc.NewClient(c.PayRpc)
		if err != nil {
			logx.Errorf("连接支付服务失败: %v", err)
		} else {
			ctx.PayRpc = pay.NewPay(payRpcClient)
			logx.Info("成功连接到支付服务")
		}
	}

	return ctx
}
