package svc

import (
	"mall/service/product/model"
	"mall/service/product/rpc/internal/config"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config         config.Config
	ProductModel   model.ProductModel
	StockProcessor *model.StockProcessor
	DB             *gorm.DB
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 创建数据库连接
	conn := sqlx.NewMysql(c.Mysql.DataSource)

	// 创建产品模型
	productModel := model.NewProductModel(conn, c.CacheRedis)

	// 创建库存处理器
	stockProcessor := model.NewStockProcessor(productModel, 10000, 10)
	// 启动库存处理器
	stockProcessor.Start()

	// 初始化GORM
	db, err := gorm.Open(mysql.Open(c.Mysql.DataSource), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	return &ServiceContext{
		Config:         c,
		ProductModel:   productModel,
		StockProcessor: stockProcessor,
		DB:             db,
	}
}
