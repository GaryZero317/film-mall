package svc

import (
	"database/sql"
	"mall/service/pay/api/internal/config"
	"mall/service/pay/model"
	"mall/service/pay/rpc/pay"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config   config.Config
	PayRpc   pay.Pay
	PayModel *model.GormPayModel
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
		PayRpc:   pay.NewPay(zrpc.MustNewClient(c.PayRpc)),
		PayModel: payModel,
	}
}
