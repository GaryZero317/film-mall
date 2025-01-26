package svc

import (
	"mall/service/cart/api/internal/config"
	"mall/service/product/rpc/product"

	"github.com/zeromicro/go-zero/zrpc"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config        config.Config
	DB            *gorm.DB
	ProductClient product.Product
}

func NewServiceContext(c config.Config) *ServiceContext {
	db, err := gorm.Open(mysql.Open(c.MySQL.DataSource), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	return &ServiceContext{
		Config:        c,
		DB:            db,
		ProductClient: product.NewProduct(zrpc.MustNewClient(c.ProductRpc)),
	}
}
