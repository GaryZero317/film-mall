package svc

import (
	"database/sql"
	"mall/service/user/api/internal/config"
	"mall/service/user/model"
	"mall/service/user/rpc/pb/user"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config     config.Config
	UserRpc    user.UserClient
	UserModel  model.UserModel
	AdminModel *model.GormAdminModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Mysql.DataSource)
	sqlDB, err := sql.Open("mysql", c.Mysql.DataSource)
	if err != nil {
		panic(err)
	}

	adminModel, err := model.NewGormAdminModel(sqlDB, c.CacheRedis)
	if err != nil {
		panic(err)
	}

	return &ServiceContext{
		Config:     c,
		UserRpc:    user.NewUserClient(zrpc.MustNewClient(c.UserRpc).Conn()),
		UserModel:  model.NewUserModel(conn, c.CacheRedis),
		AdminModel: adminModel,
	}
}
