package svc

import (
	"mall/service/order/api/internal/config"
	"mall/service/order/model"
	"mall/service/order/rpc/order"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config     config.Config
	OrderModel model.OrderModel
	OrderRpc   order.Order
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Mysql.DataSource)
	orderRpc := order.NewOrder(zrpc.MustNewClient(c.OrderRpc))

	return &ServiceContext{
		Config:     c,
		OrderModel: model.NewOrderModel(conn, c.CacheRedis),
		OrderRpc:   orderRpc,
	}
}
