package svc

import (
	"mall/service/product/api/internal/config"
	"mall/service/product/model"
	"mall/service/product/rpc/product"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config       config.Config
	ProductModel model.ProductModel
	ProductImage *model.ProductImage
	DB           *gorm.DB
	ProductRpc   product.Product
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Mysql.DataSource)

	// 初始化 GORM
	db, err := gorm.Open(mysql.Open(c.Mysql.DataSource), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	return &ServiceContext{
		Config:       c,
		ProductModel: model.NewProductModel(conn, c.CacheRedis),
		ProductImage: &model.ProductImage{},
		DB:           db,
		ProductRpc:   product.NewProduct(zrpc.MustNewClient(c.ProductRpc)),
	}
}
