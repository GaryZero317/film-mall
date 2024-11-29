package svc

import (
	"database/sql"
	"mall/service/order/model"
	"mall/service/order/rpc/internal/config"
	"mall/service/product/rpc/product"
	"mall/service/user/rpc/pb/user"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config

	OrderModel model.OrderModel

	UserRpc    user.UserClient
	ProductRpc product.Product
}

func NewServiceContext(c config.Config) *ServiceContext {
	sqlDB, err := sql.Open("mysql", c.Mysql.DataSource)
	if err != nil {
		panic(err)
	}

	orderModel, err := model.NewGormOrderModel(sqlDB, c.CacheRedis)
	if err != nil {
		panic(err)
	}

	return &ServiceContext{
		Config:     c,
		OrderModel: orderModel,
		UserRpc:    user.NewUserClient(zrpc.MustNewClient(c.UserRpc).Conn()),
		ProductRpc: product.NewProduct(zrpc.MustNewClient(c.ProductRpc)),
	}
}
