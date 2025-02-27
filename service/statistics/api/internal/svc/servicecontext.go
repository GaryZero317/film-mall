package svc

import (
	"mall/service/product/rpc/product"
	"mall/service/statistics/api/internal/config"
	"mall/service/statistics/model"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config                  config.Config
	ProductSalesDailyModel  model.ProductSalesDailyModel
	CategorySalesDailyModel model.CategorySalesDailyModel
	UserActivityLogModel    model.UserActivityLogModel
	ProductViewLogModel     model.ProductViewLogModel
	ProductRpc              product.Product
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.DB.DataSource)
	return &ServiceContext{
		Config:                  c,
		ProductSalesDailyModel:  model.NewProductSalesDailyModel(conn),
		CategorySalesDailyModel: model.NewCategorySalesDailyModel(conn),
		UserActivityLogModel:    model.NewUserActivityLogModel(conn),
		ProductViewLogModel:     model.NewProductViewLogModel(conn),
		ProductRpc:              product.NewProduct(zrpc.MustNewClient(c.ProductRpc)),
	}
}
