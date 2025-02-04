package svc

import (
	"mall/service/cart/api/internal/config"
	"mall/service/product/rpc/product"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config        config.Config
	DB            *gorm.DB
	Cache         *redis.Redis
	ProductClient product.Product
}

func NewServiceContext(c config.Config) *ServiceContext {
	db, err := gorm.Open(mysql.Open(c.MySQL.DataSource), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	rds := redis.MustNewRedis(redis.RedisConf{
		Host: c.CacheRedis[0].Host,
		Type: c.CacheRedis[0].Type,
		Pass: c.CacheRedis[0].Pass,
	})

	return &ServiceContext{
		Config:        c,
		DB:            db,
		Cache:         rds,
		ProductClient: product.NewProduct(zrpc.MustNewClient(c.ProductRpc)),
	}
}
