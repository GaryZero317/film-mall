package svc

import (
	"database/sql"
	"mall/service/order/rpc/order"
	"mall/service/pay/model"
	"mall/service/pay/rpc/internal/config"
	"mall/service/user/rpc/user"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config

	PayModel model.PayModel
	OrderRpc order.Order
	UserRpc  user.User
}

func NewServiceContext(c config.Config) *ServiceContext {
	sqlDB, err := sql.Open("mysql", c.Mysql.DataSource)
	if err != nil {
		panic(err)
	}

	payModel, err := model.NewGormPayModel(sqlDB, c.CacheRedis)
	if err != nil {
		panic(err)
	}

	return &ServiceContext{
		Config:   c,
		PayModel: payModel,
		OrderRpc: order.NewOrder(zrpc.MustNewClient(c.OrderRpc)),
		UserRpc:  user.NewUser(zrpc.MustNewClient(c.UserRpc)),
	}
}
