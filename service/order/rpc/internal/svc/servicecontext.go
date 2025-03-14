package svc

import (
	"mall/service/order/model"
	"mall/service/order/rpc/internal/config"
	"mall/service/product/rpc/pb/product"

	"github.com/go-redis/redis/v8"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config          config.Config
	OrderModel      model.OrderModel
	OrderItemModel  model.OrderItemModel
	Redis           *redis.Client
	ProductRpc      product.ProductClient
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

	return &ServiceContext{
		Config:          c,
		OrderModel:      model.NewOrderModel(conn, c.CacheRedis),
		OrderItemModel:  model.NewOrderItemModel(conn, c.CacheRedis),
		Redis:           rdb,
		ProductRpc:      product.NewProductClient(zrpc.MustNewClient(c.ProductRpc).Conn()),
		OrderLockExpiry: c.OrderLockExpiry,
	}
}
