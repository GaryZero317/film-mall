package svc

import (
	"mall/service/address/api/internal/config"
	"mall/service/address/model"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config       config.Config
	AddressModel model.AddressModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.MySQL.DataSource)
	return &ServiceContext{
		Config:       c,
		AddressModel: model.NewAddressModel(conn, c.CacheRedis),
	}
}
