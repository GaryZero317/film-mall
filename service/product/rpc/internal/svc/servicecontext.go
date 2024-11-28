package svc

import (
	"database/sql"
	"mall/service/product/model"
	"mall/service/product/rpc/internal/config"
)

type ServiceContext struct {
	Config config.Config

	ProductModel model.ProductModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	sqlDB, err := sql.Open("mysql", c.Mysql.DataSource)
	if err != nil {
		panic(err)
	}

	productModel, err := model.NewGormProductModel(sqlDB, c.CacheRedis)
	if err != nil {
		panic(err)
	}

	return &ServiceContext{
		Config:       c,
		ProductModel: productModel,
	}
}
