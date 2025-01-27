package svc

import (
	"mall/service/order/api/internal/config"
	"mall/service/order/model"
	"mall/service/order/rpc/order"
	"mall/service/product/rpc/product"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config     config.Config
	OrderRpc   order.Order
	ProductRpc product.Product
	OrderModel model.OrderModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Mysql.DataSource)
	return &ServiceContext{
		Config:     c,
		OrderRpc:   order.NewOrder(zrpc.MustNewClient(c.OrderRpc)),
		ProductRpc: product.NewProduct(zrpc.MustNewClient(c.ProductRpc)),
		OrderModel: model.NewOrderModel(conn, c.CacheRedis),
	}
}
