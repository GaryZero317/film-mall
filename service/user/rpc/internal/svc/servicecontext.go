package svc

import (
	"mall/service/user/model"
	"mall/service/user/rpc/internal/config"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config     config.Config
	UserModel  model.UserModel
	AdminModel model.AdminModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Mysql.DataSource)
	return &ServiceContext{
		Config:     c,
		UserModel:  model.NewUserModel(conn, c.CacheRedis),
		AdminModel: model.NewAdminModel(conn, c.CacheRedis),
	}
}
