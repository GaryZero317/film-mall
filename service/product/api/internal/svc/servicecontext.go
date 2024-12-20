package svc

import (
	"mall/service/product/api/internal/config"
	"mall/service/product/model"
	"mall/service/product/rpc/product"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config       config.Config
	ProductModel model.ProductModel
	ProductRpc   product.Product
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Mysql.DataSource)
	return &ServiceContext{
		Config:       c,
		ProductRpc:   product.NewProduct(zrpc.MustNewClient(c.ProductRpc)),
		ProductModel: model.NewProductModel(conn, c.CacheRedis),
	}
}
